package contract

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/auth"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/handlers"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

// Real handlers against the compose Postgres (skipped when DATABASE_URL is unset).
func TestContract_AdminCoupons(t *testing.T) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	gin.SetMode(gin.TestMode)
	jwtSvc := testJWT(t)
	r := gin.New()
	admin := r.Group("/api/v1/admin", auth.AdminMiddleware(jwtSvc))
	h := handlers.NewAdminCouponHandler(services.NewAdminCouponService(repositories.NewCouponAdminRepository(pool)))
	admin.GET("/coupons", h.List)
	admin.POST("/coupons", h.Create)
	admin.GET("/coupons/:id", h.Get)
	admin.PUT("/coupons/:id", h.Update)
	admin.DELETE("/coupons/:id", h.Delete)

	adminTok, _ := jwtSvc.SignAccessToken(uuid.New(), "admin")
	custTok, _ := jwtSvc.SignAccessToken(uuid.New(), "customer")
	call := func(tok, method, path string, body any) (int, map[string]any) {
		t.Helper()
		var buf bytes.Buffer
		if body != nil {
			_ = json.NewEncoder(&buf).Encode(body)
		}
		req := httptest.NewRequest(method, path, &buf)
		req.Header.Set("Content-Type", "application/json")
		if tok != "" {
			req.Header.Set("Authorization", "Bearer "+tok)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var out map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &out)
		return w.Code, out
	}
	expect := func(got, want int, what string) {
		t.Helper()
		if got != want {
			t.Fatalf("%s: expected %d, got %d", what, want, got)
		}
	}

	sfx := uuid.NewString()[:8]
	code := "ZZ" + strings.ToUpper(sfx)
	redeemedCode := "ZR" + strings.ToUpper(sfx)
	var userID uuid.UUID
	if e := pool.QueryRow(ctx, `INSERT INTO users (email, full_name) VALUES ($1, 'zz') RETURNING id`, "zzc-"+sfx+"@example.com").Scan(&userID); e != nil {
		t.Fatal(e)
	}
	defer func() {
		_, _ = pool.Exec(ctx, `DELETE FROM coupon_users WHERE user_id = $1`, userID)
		_, _ = pool.Exec(ctx, `DELETE FROM coupons WHERE code IN ($1, $2)`, code, redeemedCode)
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	}()

	st, _ := call(custTok, http.MethodGet, "/api/v1/admin/coupons", nil)
	expect(st, http.StatusForbidden, "customer list")

	// Create: type gates which amount is stored.
	st, c := call(adminTok, http.MethodPost, "/api/v1/admin/coupons", map[string]any{
		"name": "Zz flat", "code": "zz" + sfx, "type": "flat", "flat": 50, "percent": 10,
	})
	expect(st, http.StatusCreated, "create flat")
	mustHaveKeys(t, c, "id", "name", "code", "is_enable", "type", "flat", "percent", "coupon_expiry", "created_at", "updated_at")
	if c["code"] != code || c["flat"] != float64(50) || c["percent"] != nil || c["is_enable"] != true {
		t.Errorf("unexpected created coupon: %v", c)
	}
	id, _ := c["id"].(string)
	base := "/api/v1/admin/coupons/" + id

	st, _ = call(adminTok, http.MethodPost, "/api/v1/admin/coupons", map[string]any{"name": "Dup", "code": code, "type": "flat", "flat": 5})
	expect(st, http.StatusConflict, "duplicate code")
	st, _ = call(adminTok, http.MethodPost, "/api/v1/admin/coupons", map[string]any{"name": "Bad", "code": "BAD" + sfx, "type": "percentage", "percent": 150})
	expect(st, http.StatusUnprocessableEntity, "percent out of range")
	st, _ = call(adminTok, http.MethodPost, "/api/v1/admin/coupons", map[string]any{"name": "Bad", "code": "BAD" + sfx, "type": "flat"})
	expect(st, http.StatusUnprocessableEntity, "flat without amount")
	st, _ = call(adminTok, http.MethodPost, "/api/v1/admin/coupons", map[string]any{"name": "Bad", "code": "b a d", "type": "flat", "flat": 5})
	expect(st, http.StatusUnprocessableEntity, "invalid code")

	st, list := call(adminTok, http.MethodGet, "/api/v1/admin/coupons?limit=100", nil)
	expect(st, http.StatusOK, "list")
	mustHaveKeys(t, list, "data", "meta")
	st, got := call(adminTok, http.MethodGet, base, nil)
	expect(st, http.StatusOK, "get")
	if got["code"] != code {
		t.Errorf("unexpected get: %v", got)
	}
	st, _ = call(adminTok, http.MethodGet, "/api/v1/admin/coupons/"+uuid.NewString(), nil)
	expect(st, http.StatusNotFound, "get unknown")

	// Update: disable, then switch type to percentage (flat is cleared) with an expiry.
	st, u := call(adminTok, http.MethodPut, base, map[string]any{"is_enable": false})
	expect(st, http.StatusOK, "disable")
	if u["is_enable"] != false || u["flat"] != float64(50) {
		t.Errorf("unexpected disable result: %v", u)
	}
	st, u = call(adminTok, http.MethodPut, base, map[string]any{"type": "percentage", "percent": 15, "coupon_expiry": "2030-01-01T00:00:00Z"})
	expect(st, http.StatusOK, "switch type")
	if u["type"] != "percentage" || u["percent"] != float64(15) || u["flat"] != nil || u["coupon_expiry"] == nil {
		t.Errorf("unexpected type switch result: %v", u)
	}
	st, u = call(adminTok, http.MethodPut, base, map[string]any{"coupon_expiry": nil})
	expect(st, http.StatusOK, "clear expiry")
	if u["coupon_expiry"] != nil {
		t.Errorf("expiry should be cleared: %v", u)
	}
	st, _ = call(adminTok, http.MethodPut, base, map[string]any{"type": "flat"})
	expect(st, http.StatusUnprocessableEntity, "switching to flat needs an amount")
	st, _ = call(adminTok, http.MethodPut, "/api/v1/admin/coupons/"+uuid.NewString(), map[string]any{"is_enable": true})
	expect(st, http.StatusNotFound, "update unknown")

	// A redeemed coupon cannot be deleted; an unredeemed one can.
	st, rc := call(adminTok, http.MethodPost, "/api/v1/admin/coupons", map[string]any{"name": "Redeemed", "code": redeemedCode, "type": "flat", "flat": 5})
	expect(st, http.StatusCreated, "create redeemed")
	rid, _ := rc["id"].(string)
	if _, e := pool.Exec(ctx, `INSERT INTO coupon_users (user_id, coupon_id) VALUES ($1, $2)`, userID, rid); e != nil {
		t.Fatal(e)
	}
	st, _ = call(adminTok, http.MethodDelete, "/api/v1/admin/coupons/"+rid, nil)
	expect(st, http.StatusConflict, "delete redeemed coupon")
	st, msg := call(adminTok, http.MethodDelete, base, nil)
	expect(st, http.StatusOK, "delete")
	mustHaveKeys(t, msg, "message")
	st, _ = call(adminTok, http.MethodGet, base, nil)
	expect(st, http.StatusNotFound, "get after delete")
}
