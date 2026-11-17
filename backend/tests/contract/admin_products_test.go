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
func TestContract_AdminProducts(t *testing.T) {
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
	attrSvc := services.NewAttributeService(repositories.NewAttributeRepository(pool))
	pub := services.NewProductService(repositories.NewProductRepository(pool), repositories.NewCategoryRepository(pool), attrSvc)
	r := gin.New()
	r.GET("/api/v1/products/:id", handlers.NewProductHandler(pub).Get)
	admin := r.Group("/api/v1/admin", auth.AdminMiddleware(jwtSvc))
	h := handlers.NewAdminProductHandler(services.NewAdminProductService(repositories.NewAdminProductRepository(pool), pub))
	admin.GET("/products", h.List)
	admin.POST("/products", h.Create)
	admin.GET("/products/:id", h.Get)
	admin.PUT("/products/:id", h.Update)
	admin.DELETE("/products/:id", h.Delete)
	admin.GET("/products/:id/attributes", h.ListAttributes)
	admin.POST("/products/:id/attributes", h.AssignAttribute)
	admin.DELETE("/products/:id/attributes/:attrId", h.RemoveAttribute)

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
	var catID, attrID, valueID uuid.UUID
	if e := pool.QueryRow(ctx, `INSERT INTO categories (slug, name) VALUES ($1, 'zz') RETURNING id`, "zzp-"+sfx).Scan(&catID); e != nil {
		t.Fatal(e)
	}
	if e := pool.QueryRow(ctx, `INSERT INTO product_attributes (code, name) VALUES ($1, 'zz') RETURNING id`, "zzp_"+sfx).Scan(&attrID); e != nil {
		t.Fatal(e)
	}
	if e := pool.QueryRow(ctx, `INSERT INTO product_attribute_values (attribute_id, value, slug) VALUES ($1, 'V', 'v') RETURNING id`, attrID).Scan(&valueID); e != nil {
		t.Fatal(e)
	}
	slug := "zz-prod-" + sfx
	defer func() {
		_, _ = pool.Exec(ctx, `DELETE FROM products WHERE slug = $1`, slug)
		_, _ = pool.Exec(ctx, `DELETE FROM product_attributes WHERE id = $1`, attrID)
		_, _ = pool.Exec(ctx, `DELETE FROM categories WHERE id = $1`, catID)
	}()

	// Authorization.
	st, _ := call(custTok, http.MethodGet, "/api/v1/admin/products", nil)
	expect(st, http.StatusForbidden, "customer list")

	// POST /products
	body := map[string]any{
		"name": "Zz Product", "slug": slug, "description": "desc", "price": 499.5, "discount_percent": 10,
		"stock": 7, "is_active": true, "category_ids": []string{catID.String()},
		"image_urls": []map[string]any{{"url": "https://cdn.example.com/a.jpg", "alt_text": "front", "sort_order": 1}},
	}
	st, p := call(adminTok, http.MethodPost, "/api/v1/admin/products", body)
	expect(st, http.StatusCreated, "create")
	mustHaveKeys(t, p, "id", "name", "slug", "description", "price", "discount_percent", "discounted_price",
		"stock", "attributes", "images", "categories", "is_active", "created_at", "updated_at")
	pi, _ := p["images"].([]any)
	pc, _ := p["categories"].([]any)
	if p["stock"] != float64(7) || len(pi) != 1 || len(pc) != 1 {
		t.Errorf("unexpected created product: %v", p)
	}
	id, _ := p["id"].(string)
	base := "/api/v1/admin/products/" + id

	st, _ = call(adminTok, http.MethodPost, "/api/v1/admin/products", body)
	expect(st, http.StatusConflict, "duplicate slug")
	bad := map[string]any{"name": "", "slug": "Bad Slug", "price": -1, "category_ids": []string{catID.String()}}
	st, _ = call(adminTok, http.MethodPost, "/api/v1/admin/products", bad)
	expect(st, http.StatusUnprocessableEntity, "invalid create")

	// GET list (search) and detail
	st, list := call(adminTok, http.MethodGet, "/api/v1/admin/products?q="+slug, nil)
	expect(st, http.StatusOK, "list")
	mustHaveKeys(t, list, "data", "meta")
	if data, _ := list["data"].([]any); len(data) != 1 {
		t.Errorf("search should find exactly the new product, got %v", list["data"])
	}
	st, got := call(adminTok, http.MethodGet, base, nil)
	expect(st, http.StatusOK, "get")
	if got["name"] != "Zz Product" {
		t.Errorf("unexpected get: %v", got)
	}
	st, _ = call(adminTok, http.MethodGet, "/api/v1/admin/products/"+uuid.NewString(), nil)
	expect(st, http.StatusNotFound, "get unknown")

	// PUT: partial update replaces images and stock
	st, upd := call(adminTok, http.MethodPut, base, map[string]any{
		"name": "Renamed", "stock": 20, "image_urls": []map[string]any{{"url": "https://cdn.example.com/b.jpg", "sort_order": 1}},
	})
	expect(st, http.StatusOK, "update")
	imgs, _ := upd["images"].([]any)
	if upd["name"] != "Renamed" || upd["stock"] != float64(20) || len(imgs) != 1 || upd["slug"] != slug {
		t.Errorf("unexpected update: %v", upd)
	}

	// Attributes
	st, a := call(adminTok, http.MethodPost, base+"/attributes", map[string]any{"attribute_value_id": valueID.String()})
	expect(st, http.StatusCreated, "assign attribute")
	mustHaveKeys(t, a, "attribute_value_id", "attribute_id", "code", "attribute_name", "value", "value_slug")
	st, la := call(adminTok, http.MethodGet, base+"/attributes", nil)
	expect(st, http.StatusOK, "list attributes")
	if data, _ := la["data"].([]any); len(data) != 1 {
		t.Errorf("expected 1 attribute, got %v", la["data"])
	}
	st, _ = call(adminTok, http.MethodPost, base+"/attributes", map[string]any{"attribute_value_id": uuid.NewString()})
	expect(st, http.StatusNotFound, "assign unknown value")
	st, _ = call(adminTok, http.MethodDelete, base+"/attributes/"+valueID.String(), nil)
	expect(st, http.StatusNoContent, "remove attribute")
	st, _ = call(adminTok, http.MethodDelete, base+"/attributes/"+valueID.String(), nil)
	expect(st, http.StatusNotFound, "remove attribute twice")

	// Public endpoint serves the product while active, and 404s after the soft delete.
	st, _ = call("", http.MethodGet, "/api/v1/products/"+id, nil)
	expect(st, http.StatusOK, "public get active")
	st, msg := call(adminTok, http.MethodDelete, base, nil)
	expect(st, http.StatusOK, "soft delete")
	mustHaveKeys(t, msg, "message")
	st, _ = call("", http.MethodGet, "/api/v1/products/"+id, nil)
	expect(st, http.StatusNotFound, "public get after delete")
	st, still := call(adminTok, http.MethodGet, base, nil)
	expect(st, http.StatusOK, "admin still sees soft-deleted product")
	if still["is_active"] != false {
		t.Errorf("soft-deleted product must have is_active=false")
	}
}
