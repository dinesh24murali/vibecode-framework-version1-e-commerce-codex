package contract

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/auth"
)

// POST /api/v1/admin/auth/login → 200 AuthTokenResponse with an admin user.
func TestContract_AdminLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/admin/auth/login", func(c *gin.Context) {
		resp := stubAuthResponse()
		if u, ok := resp["user"].(gin.H); ok {
			u["role"] = "admin"
		}
		c.JSON(http.StatusOK, resp)
	})
	w := doRequest(t, r, http.MethodPost, "/api/v1/admin/auth/login", map[string]any{"email": "admin@example.com", "password": "x"})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := decodeBody(t, w)
	mustHaveKeys(t, body, "access_token", "refresh_token", "token_type", "expires_in", "user")
	if u, _ := body["user"].(map[string]any); u["role"] != "admin" {
		t.Errorf("expected admin role in login response")
	}
}

// AdminMiddleware: no/invalid token → 401, customer token → 403, admin token → 200.
func TestAdminMiddleware_RoleEnforcement(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtSvc := testJWT(t)
	r := gin.New()
	r.GET("/api/v1/admin/ping", auth.AdminMiddleware(jwtSvc), func(c *gin.Context) { c.Status(http.StatusOK) })
	adminTok, _ := jwtSvc.SignAccessToken(uuid.New(), "admin")
	custTok, _ := jwtSvc.SignAccessToken(uuid.New(), "customer")

	for name, tc := range map[string]struct {
		token string
		want  int
	}{
		"missing token":  {"", http.StatusUnauthorized},
		"garbage token":  {"abc", http.StatusUnauthorized},
		"customer token": {custTok, http.StatusForbidden},
		"admin token":    {adminTok, http.StatusOK},
	} {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/ping", nil)
		if tc.token != "" {
			req.Header.Set("Authorization", "Bearer "+tc.token)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != tc.want {
			t.Errorf("%s: expected %d, got %d", name, tc.want, w.Code)
		}
	}
}
