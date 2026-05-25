package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	otpTTL    = 10 * time.Minute
	otpDigits = 6
)

// OTPService generates, stores, and verifies one-time passwords using Redis.
type OTPService struct {
	redis *redis.Client
}

// NewOTPService returns an OTPService backed by the given Redis client.
func NewOTPService(redis *redis.Client) *OTPService {
	return &OTPService{redis: redis}
}

// GenerateAndStore creates a 6-digit OTP, stores it in Redis under the key
// "otp:<purpose>:<email>" with a 10-minute TTL, and returns the OTP string.
// The OTP is logged to stdout in development (no email transport in this phase).
func (s *OTPService) GenerateAndStore(ctx context.Context, email, purpose string) (string, error) {
	otp, err := generateOTP()
	if err != nil {
		return "", fmt.Errorf("otp: generate: %w", err)
	}
	key := otpKey(purpose, email)
	if err := s.redis.Set(ctx, key, otp, otpTTL).Err(); err != nil {
		return "", fmt.Errorf("otp: store: %w", err)
	}
	return otp, nil
}

// Verify checks that the supplied OTP matches the stored value for
// "otp:<purpose>:<email>". On success the key is deleted (single-use).
// Returns false (no error) when the OTP is simply wrong; returns an error
// only on infrastructure failure.
func (s *OTPService) Verify(ctx context.Context, email, purpose, otp string) (bool, error) {
	key := otpKey(purpose, email)
	stored, err := s.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil // expired or never issued
	}
	if err != nil {
		return false, fmt.Errorf("otp: redis get: %w", err)
	}
	if stored != otp {
		return false, nil
	}
	// Delete the key immediately — single-use.
	_ = s.redis.Del(ctx, key)
	return true, nil
}

func otpKey(purpose, email string) string {
	return fmt.Sprintf("otp:%s:%s", purpose, email)
}

func generateOTP() (string, error) {
	max := big.NewInt(1_000_000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%0*d", otpDigits, n.Int64()), nil
}
