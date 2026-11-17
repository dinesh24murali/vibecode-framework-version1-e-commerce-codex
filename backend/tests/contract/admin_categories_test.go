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
func TestContract_AdminCategories(t *testing.T) {
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
	h := handlers.NewAdminCategoryHandler(services.NewAdminCategoryService(repositories.NewCategoryAdminRepository(pool)))
	admin.GET("/categories", h.List)
	admin.POST("/categories", h.Create)
	admin.GET("/categories/:id", h.Get)
	admin.PUT("/categories/:id", h.Update)
	admin.DELETE("/categories/:id", h.Delete)

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

	slug := "zz-cat-" + uuid.NewString()[:8]
	defer func() {
		_, _ = pool.Exec(ctx, `DELETE FROM products WHERE slug LIKE $1`, slug+"%")
		_, _ = pool.Exec(ctx, `DELETE FROM categories WHERE slug LIKE $1`, slug+"%")
	}()

	st, _ := call(custTok, http.MethodGet, "/api/v1/admin/categories", nil)
	expect(st, http.StatusForbidden, "customer list")

	st, c := call(adminTok, http.MethodPost, "/api/v1/admin/categories", map[string]any{"name": "Zz Cat", "slug": slug})
	expect(st, http.StatusCreated, "create")
	mustHaveKeys(t, c, "id", "name", "slug", "image_url", "created_at", "updated_at")
	id, _ := c["id"].(string)
	base := "/api/v1/admin/categories/" + id

	st, _ = call(adminTok, http.MethodPost, "/api/v1/admin/categories", map[string]any{"name": "Dup", "slug": slug})
	expect(st, http.StatusConflict, "duplicate slug")
	st, _ = call(adminTok, http.MethodPost, "/api/v1/admin/categories", map[string]any{"name": "", "slug": "Bad Slug"})
	expect(st, http.StatusUnprocessableEntity, "invalid create")

	st, list := call(adminTok, http.MethodGet, "/api/v1/admin/categories?limit=100", nil)
	expect(st, http.StatusOK, "list")
	mustHaveKeys(t, list, "data", "meta")
	st, got := call(adminTok, http.MethodGet, base, nil)
	expect(st, http.StatusOK, "get")
	if got["slug"] != slug {
		t.Errorf("unexpected get: %v", got)
	}
	st, _ = call(adminTok, http.MethodGet, "/api/v1/admin/categories/"+uuid.NewString(), nil)
	expect(st, http.StatusNotFound, "get unknown")

	st, upd := call(adminTok, http.MethodPut, base, map[string]any{"name": "Renamed"})
	expect(st, http.StatusOK, "update")
	if upd["name"] != "Renamed" || upd["slug"] != slug {
		t.Errorf("unexpected update: %v", upd)
	}

	// A category that still has a product (even an inactive one) cannot be deleted.
	if _, e := pool.Exec(ctx, `INSERT INTO products (category_id, sku, slug, title, price_inr) VALUES ($1, $2, $2, 'zz', 1)`, id, slug+"-p"); e != nil {
		t.Fatal(e)
	}
	st, _ = call(adminTok, http.MethodDelete, base, nil)
	expect(st, http.StatusConflict, "delete category with products")
	if _, e := pool.Exec(ctx, `DELETE FROM products WHERE slug = $1`, slug+"-p"); e != nil {
		t.Fatal(e)
	}
	st, msg := call(adminTok, http.MethodDelete, base, nil)
	expect(st, http.StatusOK, "delete")
	mustHaveKeys(t, msg, "message")
	st, _ = call(adminTok, http.MethodGet, base, nil)
	expect(st, http.StatusNotFound, "get after delete")
}
