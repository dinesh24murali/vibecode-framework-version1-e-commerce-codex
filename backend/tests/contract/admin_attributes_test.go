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

// Runs the real admin attribute handlers against the compose Postgres
// (skipped when DATABASE_URL is unset) and checks every endpoint against the spec.
func TestContract_AdminAttributes(t *testing.T) {
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
	h := handlers.NewAdminAttributeHandler(services.NewAttributeAdminService(repositories.NewAttributeAdminRepository(pool)))
	admin.GET("/attributes", h.ListDefinitions)
	admin.POST("/attributes", h.CreateDefinition)
	admin.GET("/attributes/:attrDefId", h.GetDefinition)
	admin.PUT("/attributes/:attrDefId", h.UpdateDefinition)
	admin.DELETE("/attributes/:attrDefId", h.DeleteDefinition)
	admin.GET("/attributes/:attrDefId/values", h.ListValues)
	admin.POST("/attributes/:attrDefId/values", h.CreateValue)
	admin.PUT("/attributes/:attrDefId/values/:valueId", h.UpdateValue)
	admin.DELETE("/attributes/:attrDefId/values/:valueId", h.DeleteValue)

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

	code := "zz_" + uuid.NewString()[:8]
	defer func() { _, _ = pool.Exec(ctx, `DELETE FROM product_attributes WHERE code = $1`, code) }()

	// Authorization: customer token is 403 on every admin attribute route.
	st, _ := call(custTok, http.MethodGet, "/api/v1/admin/attributes", nil)
	expect(st, http.StatusForbidden, "customer list")
	st, _ = call("", http.MethodGet, "/api/v1/admin/attributes", nil)
	expect(st, http.StatusUnauthorized, "anonymous list")

	// POST /attributes
	st, def := call(adminTok, http.MethodPost, "/api/v1/admin/attributes", map[string]any{"code": code, "name": "Zz Test", "sort_order": 99})
	expect(st, http.StatusCreated, "create definition")
	mustHaveKeys(t, def, "id", "code", "name", "is_filterable", "sort_order", "created_at", "updated_at")
	if def["is_filterable"] != true {
		t.Errorf("is_filterable should default to true")
	}
	id, _ := def["id"].(string)
	base := "/api/v1/admin/attributes/" + id

	st, _ = call(adminTok, http.MethodPost, "/api/v1/admin/attributes", map[string]any{"code": code, "name": "Dup"})
	expect(st, http.StatusConflict, "duplicate code")
	st, _ = call(adminTok, http.MethodPost, "/api/v1/admin/attributes", map[string]any{"code": "Bad Code", "name": "X"})
	expect(st, http.StatusUnprocessableEntity, "invalid code")

	// GET /attributes (paginated) and /attributes/:id
	st, list := call(adminTok, http.MethodGet, "/api/v1/admin/attributes?page=1&limit=100", nil)
	expect(st, http.StatusOK, "list definitions")
	mustHaveKeys(t, list, "data", "meta")
	meta, _ := list["meta"].(map[string]any)
	mustHaveKeys(t, meta, "total", "page", "limit", "total_pages")
	st, _ = call(adminTok, http.MethodGet, "/api/v1/admin/attributes?limit=0", nil)
	expect(st, http.StatusBadRequest, "bad pagination")
	st, got := call(adminTok, http.MethodGet, base, nil)
	expect(st, http.StatusOK, "get definition")
	mustHaveKeys(t, got, "id", "code", "name")
	st, _ = call(adminTok, http.MethodGet, "/api/v1/admin/attributes/"+uuid.NewString(), nil)
	expect(st, http.StatusNotFound, "get unknown definition")

	// PUT /attributes/:id (rename)
	st, upd := call(adminTok, http.MethodPut, base, map[string]any{"name": "Renamed"})
	expect(st, http.StatusOK, "update definition")
	if upd["name"] != "Renamed" || upd["code"] != code {
		t.Errorf("unexpected update result %v", upd)
	}

	// Values
	st, v1 := call(adminTok, http.MethodPost, base+"/values", map[string]any{"value": "Beginner Level", "sort_order": 2})
	expect(st, http.StatusCreated, "create value")
	mustHaveKeys(t, v1, "id", "attribute_id", "value", "value_slug", "sort_order", "created_at", "updated_at")
	if v1["value_slug"] != "beginner-level" {
		t.Errorf("slug should be generated, got %v", v1["value_slug"])
	}
	st, _ = call(adminTok, http.MethodPost, base+"/values", map[string]any{"value": "Beginner Level"})
	expect(st, http.StatusConflict, "duplicate slug")
	st, _ = call(adminTok, http.MethodPost, base+"/values", map[string]any{"value": "  "})
	expect(st, http.StatusUnprocessableEntity, "empty value")
	vid, _ := v1["id"].(string)

	st, vals := call(adminTok, http.MethodGet, base+"/values", nil)
	expect(st, http.StatusOK, "list values")
	if data, _ := vals["data"].([]any); len(data) != 1 {
		t.Errorf("expected 1 value, got %v", vals["data"])
	}
	st, vupd := call(adminTok, http.MethodPut, base+"/values/"+vid, map[string]any{"value": "Beginner"})
	expect(st, http.StatusOK, "update value")
	if vupd["value"] != "Beginner" {
		t.Errorf("value not renamed: %v", vupd)
	}
	st, _ = call(adminTok, http.MethodPut, base+"/values/"+uuid.NewString(), map[string]any{"value": "X"})
	expect(st, http.StatusNotFound, "update unknown value")

	// Delete value assigned to a product → 409.
	var catID, productID uuid.UUID
	sfx := uuid.NewString()[:8]
	if e := pool.QueryRow(ctx, `INSERT INTO categories (slug, name) VALUES ($1, 'zz') RETURNING id`, "zza-"+sfx).Scan(&catID); e != nil {
		t.Fatal(e)
	}
	defer func() { _, _ = pool.Exec(ctx, `DELETE FROM categories WHERE id = $1`, catID) }()
	if e := pool.QueryRow(ctx, `INSERT INTO products (category_id, sku, slug, title, price_inr) VALUES ($1, $2, $2, 'zz', 100) RETURNING id`, catID, "zza-"+sfx).Scan(&productID); e != nil {
		t.Fatal(e)
	}
	defer func() { _, _ = pool.Exec(ctx, `DELETE FROM products WHERE id = $1`, productID) }()
	if _, e := pool.Exec(ctx, `INSERT INTO product_attribute_assignments (product_id, attribute_value_id) VALUES ($1, $2)`, productID, vid); e != nil {
		t.Fatal(e)
	}
	st, _ = call(adminTok, http.MethodDelete, base+"/values/"+vid, nil)
	expect(st, http.StatusConflict, "delete assigned value")

	// Once unassigned, the value can be deleted (204), then the definition (204).
	if _, e := pool.Exec(ctx, `DELETE FROM product_attribute_assignments WHERE product_id = $1`, productID); e != nil {
		t.Fatal(e)
	}
	st, _ = call(adminTok, http.MethodDelete, base+"/values/"+vid, nil)
	expect(st, http.StatusNoContent, "delete value")
	st, _ = call(adminTok, http.MethodDelete, base+"/values/"+vid, nil)
	expect(st, http.StatusNotFound, "delete value twice")
	st, _ = call(adminTok, http.MethodDelete, base, nil)
	expect(st, http.StatusNoContent, "delete definition")
	st, _ = call(adminTok, http.MethodGet, base, nil)
	expect(st, http.StatusNotFound, "get deleted definition")
}
