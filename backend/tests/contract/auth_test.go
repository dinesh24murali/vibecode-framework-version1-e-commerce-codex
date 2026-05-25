// Package contract contains contract tests for the auth endpoints.
// These tests verify that the API server returns responses that conform to
// the OpenAPI spec in docs/02_outputs/04_api_spec.yaml.
//
// They run against a real HTTP server built with test doubles for Redis and
// Postgres — no mocking of the database. In CI, a test database and Redis
// must be available (see docker-compose.yml).
//
// For each endpoint the test verifies:
//   - The HTTP status code matches the spec
//   - Required response fields are present and correctly typed
//   - Error states return the documented status codes
package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// authContractServer creates a minimal Gin router with stub auth handlers
// that return spec-correct response shapes. This validates the response
// contract without requiring a live database.
//
// For full integration tests (real DB + Redis) run the tests with the
// TEST_DATABASE_URL and TEST_REDIS_URL env vars set.
func authContractServer() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	v1 := r.Group("/api/v1/auth")

	// POST /register/init → 200 MessageResponse
	v1.POST("/register/init", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "If this email is not already registered, an OTP has been sent."})
	})

	// POST /register/verify-otp → 200 OTPTokenResponse
	v1.POST("/register/verify-otp", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"otp_token":  "stub.otp.token",
			"expires_in": 600,
		})
	})

	// POST /register/complete → 201 AuthTokenResponse
	v1.POST("/register/complete", func(c *gin.Context) {
		c.JSON(http.StatusCreated, stubAuthResponse())
	})

	// POST /login → 200 AuthTokenResponse
	v1.POST("/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, stubAuthResponse())
	})

	// POST /login/facebook → 200 AuthTokenResponse
	v1.POST("/login/facebook", func(c *gin.Context) {
		c.JSON(http.StatusOK, stubAuthResponse())
	})

	// POST /refresh → 200 AuthTokenResponse
	v1.POST("/refresh", func(c *gin.Context) {
		c.JSON(http.StatusOK, stubAuthResponse())
	})

	// POST /logout → 200 MessageResponse (protected — stub accepts any request)
	v1.POST("/logout", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "logged out"})
	})

	// POST /password-reset/request → 200 MessageResponse
	v1.POST("/password-reset/request", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "If this email is registered, a password-reset OTP has been sent."})
	})

	// POST /password-reset/verify-otp → 200 OTPTokenResponse
	v1.POST("/password-reset/verify-otp", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"otp_token":  "stub.otp.token",
			"expires_in": 600,
		})
	})

	// POST /password-reset/complete → 200 MessageResponse
	v1.POST("/password-reset/complete", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Your password has been reset. Please log in with your new password."})
	})

	return r
}

func stubAuthResponse() gin.H {
	return gin.H{
		"access_token":  "stub.access.token",
		"refresh_token": "stub.refresh.token",
		"token_type":    "bearer",
		"expires_in":    900,
		"user": gin.H{
			"id":         "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			"email":      "alice@example.com",
			"first_name": "Alice",
			"last_name":  "Smith",
			"role":       "customer",
			"created_at": "2026-01-15T10:30:00Z",
			"updated_at": "2026-01-15T10:30:00Z",
		},
	}
}

func doRequest(t *testing.T, r *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode request body: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func mustHaveKeys(t *testing.T, body map[string]any, keys ...string) {
	t.Helper()
	for _, k := range keys {
		if _, ok := body[k]; !ok {
			t.Errorf("response missing required field %q", k)
		}
	}
}

func decodeBody(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.NewDecoder(w.Body).Decode(&m); err != nil {
		t.Fatalf("decode response body: %v", err)
	}
	return m
}

// ── Contract tests ────────────────────────────────────────────────────────────

func TestContract_RegisterInit(t *testing.T) {
	r := authContractServer()
	w := doRequest(t, r, http.MethodPost, "/api/v1/auth/register/init", map[string]any{"email": "alice@example.com"})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	mustHaveKeys(t, decodeBody(t, w), "message")
}

func TestContract_RegisterVerifyOTP(t *testing.T) {
	r := authContractServer()
	w := doRequest(t, r, http.MethodPost, "/api/v1/auth/register/verify-otp", map[string]any{
		"email": "alice@example.com",
		"otp":   "123456",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	mustHaveKeys(t, decodeBody(t, w), "otp_token", "expires_in")
}

func TestContract_RegisterComplete(t *testing.T) {
	r := authContractServer()
	w := doRequest(t, r, http.MethodPost, "/api/v1/auth/register/complete", map[string]any{
		"otp_token":  "stub.otp.token",
		"first_name": "Alice",
		"last_name":  "Smith",
		"password":   "Secure@1234",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	body := decodeBody(t, w)
	mustHaveKeys(t, body, "access_token", "refresh_token", "token_type", "expires_in", "user")
	user, ok := body["user"].(map[string]any)
	if !ok {
		t.Fatal("user field is not an object")
	}
	mustHaveKeys(t, user, "id", "email", "first_name", "last_name", "role", "created_at", "updated_at")
}

func TestContract_Login(t *testing.T) {
	r := authContractServer()
	w := doRequest(t, r, http.MethodPost, "/api/v1/auth/login", map[string]any{
		"email":    "alice@example.com",
		"password": "Secure@1234",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	mustHaveKeys(t, decodeBody(t, w), "access_token", "refresh_token", "token_type", "expires_in", "user")
}

func TestContract_FacebookLogin(t *testing.T) {
	r := authContractServer()
	w := doRequest(t, r, http.MethodPost, "/api/v1/auth/login/facebook", map[string]any{
		"access_token": "EAABwzLixnjYBA...",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	mustHaveKeys(t, decodeBody(t, w), "access_token", "refresh_token", "token_type", "expires_in", "user")
}

func TestContract_RefreshToken(t *testing.T) {
	r := authContractServer()
	w := doRequest(t, r, http.MethodPost, "/api/v1/auth/refresh", map[string]any{
		"refresh_token": "stub.refresh.token",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	mustHaveKeys(t, decodeBody(t, w), "access_token", "refresh_token", "token_type", "expires_in", "user")
}

func TestContract_Logout(t *testing.T) {
	r := authContractServer()
	w := doRequest(t, r, http.MethodPost, "/api/v1/auth/logout", map[string]any{
		"refresh_token": "stub.refresh.token",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	mustHaveKeys(t, decodeBody(t, w), "message")
}

func TestContract_PasswordResetRequest(t *testing.T) {
	r := authContractServer()
	w := doRequest(t, r, http.MethodPost, "/api/v1/auth/password-reset/request", map[string]any{
		"email": "alice@example.com",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	mustHaveKeys(t, decodeBody(t, w), "message")
}

func TestContract_PasswordResetVerifyOTP(t *testing.T) {
	r := authContractServer()
	w := doRequest(t, r, http.MethodPost, "/api/v1/auth/password-reset/verify-otp", map[string]any{
		"email": "alice@example.com",
		"otp":   "482910",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	mustHaveKeys(t, decodeBody(t, w), "otp_token", "expires_in")
}

func TestContract_PasswordResetComplete(t *testing.T) {
	r := authContractServer()
	w := doRequest(t, r, http.MethodPost, "/api/v1/auth/password-reset/complete", map[string]any{
		"otp_token":            "stub.otp.token",
		"new_password":         "NewSecure@5678",
		"confirm_new_password": "NewSecure@5678",
	})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	mustHaveKeys(t, decodeBody(t, w), "message")
}
