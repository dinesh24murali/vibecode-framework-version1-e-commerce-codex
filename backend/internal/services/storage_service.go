package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// Image upload limits from the API spec (ImagePresignRequest).
const (
	MaxImageBytes = 5 * 1024 * 1024
	// PresignTTL is how long a presigned upload URL stays valid (spec: 5 minutes).
	PresignTTL = 5 * time.Minute
)

// Sentinel errors for storage.
var (
	ErrStorageNotConfigured = errors.New("image storage is not configured")
	ErrInvalidImage         = errors.New("unsupported image type or size")
)

var imageExt = map[string]string{"image/jpeg": "jpg", "image/png": "png", "image/webp": "webp"}

// PresignedUpload is the result of a presign request.
type PresignedUpload struct {
	ExpiresAt    time.Time
	PresignedURL string
	PublicURL    string
	Key          string
}

// StorageService issues presigned URLs for product images.
type StorageService interface {
	// PresignUpload returns a presigned PUT URL for a new object plus its public URL.
	PresignUpload(ctx context.Context, contentType string, size int64) (*PresignedUpload, error)
	// PresignRead returns a presigned GET URL for an existing object key.
	PresignRead(ctx context.Context, key string, ttl time.Duration) (string, error)
}

// S3Adapter implements StorageService on AWS S3. Credentials come from the
// standard AWS chain (AWS_ACCESS_KEY_ID / AWS_SECRET_ACCESS_KEY / role).
type S3Adapter struct {
	presign *s3.PresignClient
	now     func() time.Time
	bucket  string
	baseURL string
}

// NewS3Adapter builds an adapter for bucket in region. publicBaseURL is the CDN
// (or bucket) origin that serves uploaded objects, e.g. https://cdn.example.com.
func NewS3Adapter(ctx context.Context, bucket, region, publicBaseURL string) (*S3Adapter, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("storage: load aws config: %w", err)
	}
	return newS3Adapter(s3.NewFromConfig(cfg), bucket, publicBaseURL), nil
}

func newS3Adapter(client *s3.Client, bucket, publicBaseURL string) *S3Adapter {
	if publicBaseURL == "" {
		publicBaseURL = fmt.Sprintf("https://%s.s3.amazonaws.com", bucket)
	}
	return &S3Adapter{
		presign: s3.NewPresignClient(client), bucket: bucket,
		baseURL: strings.TrimRight(publicBaseURL, "/"), now: time.Now,
	}
}

// PresignUpload implements StorageService. The content type and length are part
// of the signature, so the client must upload exactly what it declared.
func (a *S3Adapter) PresignUpload(ctx context.Context, contentType string, size int64) (*PresignedUpload, error) {
	ext, ok := imageExt[contentType]
	if !ok || size < 1 || size > MaxImageBytes {
		return nil, ErrInvalidImage
	}
	key := fmt.Sprintf("products/%s.%s", uuid.NewString(), ext)
	req, err := a.presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(a.bucket), Key: aws.String(key),
		ContentType: aws.String(contentType), ContentLength: aws.Int64(size),
	}, s3.WithPresignExpires(PresignTTL))
	if err != nil {
		return nil, fmt.Errorf("storage: presign put: %w", err)
	}
	return &PresignedUpload{
		PresignedURL: req.URL, PublicURL: a.baseURL + "/" + key, Key: key,
		ExpiresAt: a.now().Add(PresignTTL).UTC(),
	}, nil
}

// PresignRead implements StorageService.
func (a *S3Adapter) PresignRead(ctx context.Context, key string, ttl time.Duration) (string, error) {
	req, err := a.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(a.bucket), Key: aws.String(key),
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", fmt.Errorf("storage: presign get: %w", err)
	}
	return req.URL, nil
}

// DisabledStorage is used when no bucket is configured; every call fails with
// ErrStorageNotConfigured so the API still starts in local development.
type DisabledStorage struct{}

// PresignUpload implements StorageService.
func (DisabledStorage) PresignUpload(context.Context, string, int64) (*PresignedUpload, error) {
	return nil, ErrStorageNotConfigured
}

// PresignRead implements StorageService.
func (DisabledStorage) PresignRead(context.Context, string, time.Duration) (string, error) {
	return "", ErrStorageNotConfigured
}
