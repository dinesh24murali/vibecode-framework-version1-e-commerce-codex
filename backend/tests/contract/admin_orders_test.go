package contract

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
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
func TestContract_AdminOrders(t *testing.T) {
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
	h := handlers.NewAdminOrderHandler(services.NewAdminOrderService(
		repositories.NewAdminOrderRepository(pool), repositories.NewUserRepository(pool)))
	admin.GET("/orders", h.List)
	admin.GET("/orders/:id", h.Get)
	admin.PUT("/orders/:id/status", h.UpdateStatus)

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
	var userID, paidID, pendingID uuid.UUID
	if e := pool.QueryRow(ctx, `INSERT INTO users (email, full_name) VALUES ($1, 'Zz Customer') RETURNING id`, "zzo-"+sfx+"@example.com").Scan(&userID); e != nil {
		t.Fatal(e)
	}
	defer func() {
		_, _ = pool.Exec(ctx, `DELETE FROM orders WHERE user_id = $1`, userID)
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	}()
	insert := func(num, status string) uuid.UUID {
		var id uuid.UUID
		if e := pool.QueryRow(ctx, `INSERT INTO orders (user_id, order_number, status, subtotal_inr, total_inr) VALUES ($1, $2, $3, 100, 100) RETURNING id`,
			userID, num, status).Scan(&id); e != nil {
			t.Fatal(e)
		}
		return id
	}
	paidID = insert("ZZ-"+sfx+"-1", "paid")
	pendingID = insert("ZZ-"+sfx+"-2", "pending_payment")

	st, _ := call(custTok, http.MethodGet, "/api/v1/admin/orders", nil)
	expect(st, http.StatusForbidden, "customer list")

	// List (all customers, status filter, pagination)
	st, list := call(adminTok, http.MethodGet, "/api/v1/admin/orders?status=open&limit=100", nil)
	expect(st, http.StatusOK, "list open")
	mustHaveKeys(t, list, "data", "meta")
	meta, _ := list["meta"].(map[string]any)
	mustHaveKeys(t, meta, "total", "page", "limit", "total_pages")
	st, _ = call(adminTok, http.MethodGet, "/api/v1/admin/orders?status=bogus", nil)
	expect(st, http.StatusBadRequest, "bad status filter")
	st, _ = call(adminTok, http.MethodGet, "/api/v1/admin/orders?limit=0", nil)
	expect(st, http.StatusBadRequest, "bad pagination")

	// Get
	st, o := call(adminTok, http.MethodGet, "/api/v1/admin/orders/"+paidID.String(), nil)
	expect(st, http.StatusOK, "get")
	mustHaveKeys(t, o, "id", "order_number", "status", "items", "shipping_address", "subtotal", "coupon_discount",
		"taxable_amount", "gst_amount", "total", "created_at", "updated_at", "customer_name", "customer_email")
	if o["status"] != "open" || o["customer_name"] != "Zz Customer" {
		t.Errorf("unexpected order: %v", o)
	}
	st, _ = call(adminTok, http.MethodGet, "/api/v1/admin/orders/"+uuid.NewString(), nil)
	expect(st, http.StatusNotFound, "get unknown")

	// Status transitions
	put := func(id uuid.UUID, s string) (int, map[string]any) {
		return call(adminTok, http.MethodPut, "/api/v1/admin/orders/"+id.String()+"/status", map[string]any{"status": s})
	}
	st, _ = put(pendingID, "fulfilled")
	expect(st, http.StatusUnprocessableEntity, "unpaid order cannot be fulfilled")
	st, _ = put(paidID, "open")
	expect(st, http.StatusUnprocessableEntity, "cannot set open")
	st, f := put(paidID, "fulfilled")
	expect(st, http.StatusOK, "fulfil")
	if f["status"] != "fulfilled" {
		t.Errorf("expected fulfilled, got %v", f["status"])
	}
	st, _ = put(paidID, "cancelled")
	expect(st, http.StatusUnprocessableEntity, "fulfilled order is final")
	st, c := put(pendingID, "cancelled")
	expect(st, http.StatusOK, "cancel")
	if c["status"] != "cancelled" {
		t.Errorf("expected cancelled, got %v", c["status"])
	}
	st, _ = put(uuid.New(), "cancelled")
	expect(st, http.StatusNotFound, "status of unknown order")

	// Filters see the moved orders.
	st, ful := call(adminTok, http.MethodGet, "/api/v1/admin/orders?status=fulfilled&limit=100", nil)
	expect(st, http.StatusOK, "list fulfilled")
	found := false
	data, _ := ful["data"].([]any)
	for _, d := range data {
		if m, _ := d.(map[string]any); m["id"] == paidID.String() {
			found = true
		}
	}
	if !found {
		t.Errorf("fulfilled filter should include the fulfilled order")
	}
}
