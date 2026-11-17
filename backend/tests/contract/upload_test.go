package contract

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/auth"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/handlers"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

func testJWT(t *testing.T) *auth.JWTService {
	t.Helper()
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	privDER, _ := x509.MarshalPKCS8PrivateKey(priv)
	pubDER, _ := x509.MarshalPKIXPublicKey(pub)
	svc, err := auth.NewJWTService(
		string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privDER})),
		string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER})),
	)
	if err != nil {
		t.Fatal(err)
	}
	return svc
}

func uploadServer(t *testing.T, storage services.StorageService) (*gin.Engine, *auth.JWTService) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	jwtSvc := testJWT(t)
	r := gin.New()
	admin := r.Group("/api/v1/admin", auth.Middleware(jwtSvc), auth.RequireRole("admin"))
	admin.POST("/upload/image/presign", handlers.NewUploadHandler(storage).Presign)
	return r, jwtSvc
}

func presign(r *gin.Engine, token string, body any) (int, map[string]any) {
	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/upload/image/presign", &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var out map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &out)
	return w.Code, out
}

func TestUpload_Presign(t *testing.T) {
	// Static fake credentials: presigning is a local computation, no AWS call is made.
	t.Setenv("AWS_ACCESS_KEY_ID", "AKIDEXAMPLE")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
	s3a, err := services.NewS3Adapter(context.Background(), "test-bucket", "ap-south-1", "https://cdn.example.com/")
	if err != nil {
		t.Fatal(err)
	}
	r, jwtSvc := uploadServer(t, s3a)
	adminTok, _ := jwtSvc.SignAccessToken(uuid.New(), "admin")
	custTok, _ := jwtSvc.SignAccessToken(uuid.New(), "customer")
	valid := gin.H{"filename": "cover.png", "content_type": "image/png", "file_size_bytes": 1024}

	code, body := presign(r, adminTok, valid)
	if code != http.StatusOK {
		t.Fatalf("status %d %v", code, body)
	}
	requireKeys(t, body, "presigned_url", "public_url", "expires_at")

	u, err := url.Parse(body["presigned_url"].(string)) //nolint:errcheck // shape asserted above
	if err != nil || !strings.Contains(u.Host, "test-bucket") || !strings.HasPrefix(u.Path, "/products/") || !strings.HasSuffix(u.Path, ".png") {
		t.Fatalf("presigned url %v err %v", u, err)
	}
	if u.Query().Get("X-Amz-Expires") != "300" || u.Query().Get("X-Amz-Signature") == "" {
		t.Fatalf("must be signed and expire in 5 minutes: %v", u.Query())
	}
	pub, _ := body["public_url"].(string)
	if !strings.HasPrefix(pub, "https://cdn.example.com/products/") || !strings.HasSuffix(pub, ".png") {
		t.Fatalf("public url %q", pub)
	}
	exp, err := time.Parse(time.RFC3339, body["expires_at"].(string)) //nolint:errcheck // shape asserted above
	if err != nil || time.Until(exp) < 4*time.Minute || time.Until(exp) > 6*time.Minute {
		t.Fatalf("expires_at %v err %v", exp, err)
	}

	// Validation: type and size limits from the spec.
	for _, bad := range []gin.H{
		{"filename": "a.gif", "content_type": "image/gif", "file_size_bytes": 10},
		{"filename": "a.png", "content_type": "image/png", "file_size_bytes": 5*1024*1024 + 1},
		{"filename": "a.png", "content_type": "image/png"},
	} {
		if code, body := presign(r, adminTok, bad); code != http.StatusBadRequest {
			t.Errorf("%v: status %d %v", bad, code, body)
		}
	}

	// Authorization: no token -> 401, non-admin -> 403.
	if code, _ := presign(r, "", valid); code != http.StatusUnauthorized {
		t.Errorf("no token: %d", code)
	}
	if code, _ := presign(r, custTok, valid); code != http.StatusForbidden {
		t.Errorf("customer: %d", code)
	}
}

func TestUpload_NotConfigured(t *testing.T) {
	r, jwtSvc := uploadServer(t, services.DisabledStorage{})
	tok, _ := jwtSvc.SignAccessToken(uuid.New(), "admin")
	code, _ := presign(r, tok, gin.H{"filename": "a.png", "content_type": "image/png", "file_size_bytes": 10})
	if code != http.StatusServiceUnavailable {
		t.Fatalf("status %d", code)
	}
}
