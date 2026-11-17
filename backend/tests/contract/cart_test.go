package contract

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/auth"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/handlers"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

type fakeCart struct{ rows []dbsqlc.Cart }

func (f *fakeCart) ListByUser(_ context.Context, uid uuid.UUID) ([]dbsqlc.Cart, error) {
	var out []dbsqlc.Cart
	for _, r := range f.rows {
		if r.UserID == uid {
			out = append(out, r)
		}
	}
	return out, nil
}

func (f *fakeCart) GetByProduct(_ context.Context, uid, pid uuid.UUID) (*dbsqlc.Cart, error) {
	for i := range f.rows {
		if f.rows[i].UserID == uid && f.rows[i].ProductID == pid {
			return &f.rows[i], nil
		}
	}
	return nil, sql.ErrNoRows
}

func (f *fakeCart) Upsert(_ context.Context, uid, pid uuid.UUID, q int32) error {
	for i := range f.rows {
		if f.rows[i].UserID == uid && f.rows[i].ProductID == pid {
			f.rows[i].Quantity = q
			return nil
		}
	}
	f.rows = append(f.rows, dbsqlc.Cart{ID: uuid.New(), UserID: uid, ProductID: pid, Quantity: q, UpdatedAt: time.Now()})
	return nil
}

func (f *fakeCart) SetQuantity(_ context.Context, id, uid uuid.UUID, q int32) error {
	for i := range f.rows {
		if f.rows[i].ID == id && f.rows[i].UserID == uid {
			f.rows[i].Quantity = q
			return nil
		}
	}
	return sql.ErrNoRows
}

func (f *fakeCart) Delete(_ context.Context, id, uid uuid.UUID) (int64, error) {
	for i := range f.rows {
		if f.rows[i].ID == id && f.rows[i].UserID == uid {
			f.rows = append(f.rows[:i], f.rows[i+1:]...)
			return 1, nil
		}
	}
	return 0, nil
}

func (f *fakeCart) Clear(_ context.Context, uid uuid.UUID) error {
	kept := f.rows[:0]
	for _, r := range f.rows {
		if r.UserID != uid {
			kept = append(kept, r)
		}
	}
	f.rows = kept
	return nil
}

type cartEnv struct {
	r       *gin.Engine
	product uuid.UUID
}

func newCartEnv() *cartEnv {
	gin.SetMode(gin.TestMode)
	now := time.Now()
	cat := dbsqlc.Category{ID: uuid.New(), Slug: "journals", Name: "Journals", CreatedAt: now, UpdatedAt: now}
	fp := &fakeProducts{product: dbsqlc.GetActiveProductByIDRow{
		ID: uuid.New(), CategoryID: cat.ID, Slug: "mandala-journal", Title: "Mandala Journal", PriceInr: "500.00",
		DiscountPercent: sql.NullInt32{Int32: 10, Valid: true}, IsActive: true, CreatedAt: now, UpdatedAt: now,
	}}
	attrSvc := services.NewAttributeService(&fakeAttrs{values: map[uuid.UUID][]dbsqlc.ProductAttributeValue{}})
	products := services.NewProductService(fp, &fakeCategories{rows: []dbsqlc.Category{cat}}, attrSvc)
	h := handlers.NewCartHandler(services.NewCartService(&fakeCart{}, products))

	me := uuid.New()
	r := gin.New()
	g := r.Group("/api/v1", func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "unauthorized"})
			return
		}
		c.Set(auth.ContextUserID, me)
	})
	g.GET("/cart", h.Get)
	g.DELETE("/cart", h.Clear)
	g.POST("/cart/items", h.AddItem)
	g.PUT("/cart/items/:id", h.UpdateItem)
	g.DELETE("/cart/items/:id", h.RemoveItem)
	return &cartEnv{r: r, product: fp.product.ID}
}

func (e *cartEnv) do(method, path string, body any, authed bool) (int, map[string]any) {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if authed {
		req.Header.Set("Authorization", "Bearer x")
	}
	w := httptest.NewRecorder()
	e.r.ServeHTTP(w, req)
	var out map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &out)
	return w.Code, out
}

func TestCart_RequiresAuth(t *testing.T) {
	e := newCartEnv()
	id := uuid.NewString()
	for _, c := range []struct{ m, p string }{
		{"GET", "/api/v1/cart"}, {"DELETE", "/api/v1/cart"}, {"POST", "/api/v1/cart/items"},
		{"PUT", "/api/v1/cart/items/" + id}, {"DELETE", "/api/v1/cart/items/" + id},
	} {
		if code, _ := e.do(c.m, c.p, map[string]any{}, false); code != http.StatusUnauthorized {
			t.Fatalf("%s %s: status %d", c.m, c.p, code)
		}
	}
}

func TestCart_EmptyShape(t *testing.T) {
	e := newCartEnv()
	code, body := e.do("GET", "/api/v1/cart", nil, true)
	requireStatus(t, code, http.StatusOK, body)
	requireKeys(t, body, "id", "items", "item_count", "subtotal", "updated_at")
	if len(arr(body["items"])) != 0 || body["subtotal"] != float64(0) {
		t.Fatalf("expected empty cart: %v", body)
	}
}

func TestCart_AddIncrementsAndComputesTotal(t *testing.T) {
	e := newCartEnv()
	add := map[string]any{"product_id": e.product, "quantity": 2}
	code, body := e.do("POST", "/api/v1/cart/items", add, true)
	requireStatus(t, code, http.StatusCreated, body)
	code, body = e.do("POST", "/api/v1/cart/items", add, true)
	requireStatus(t, code, http.StatusCreated, body)

	items := arr(body["items"])
	if len(items) != 1 || body["item_count"] != float64(1) {
		t.Fatalf("duplicate line created: %v", body)
	}
	item := obj(items[0])
	requireKeys(t, item, "id", "product", "quantity", "unit_price", "line_total")
	// 500 with 10% discount = 450 each, 4 units.
	if item["quantity"] != float64(4) || item["unit_price"] != float64(450) || item["line_total"] != float64(1800) || body["subtotal"] != float64(1800) {
		t.Fatalf("server-side totals wrong: %v", body)
	}
}

func TestCart_AddValidation(t *testing.T) {
	e := newCartEnv()
	if code, _ := e.do("POST", "/api/v1/cart/items", map[string]any{"product_id": e.product, "quantity": 0}, true); code != http.StatusBadRequest {
		t.Fatalf("qty 0: %d", code)
	}
	if code, _ := e.do("POST", "/api/v1/cart/items", map[string]any{"product_id": uuid.New(), "quantity": 1}, true); code != http.StatusNotFound {
		t.Fatalf("unknown product: %d", code)
	}
	if code, _ := e.do("POST", "/api/v1/cart/items", map[string]any{"quantity": 1}, true); code != http.StatusBadRequest {
		t.Fatalf("missing product: %d", code)
	}
}

func TestCart_UpdateAndRemove(t *testing.T) {
	e := newCartEnv()
	_, body := e.do("POST", "/api/v1/cart/items", map[string]any{"product_id": e.product, "quantity": 1}, true)
	itemID, _ := obj(arr(body["items"])[0])["id"].(string)
	path := "/api/v1/cart/items/" + itemID

	code, body := e.do("PUT", path, map[string]any{"quantity": 3}, true)
	requireStatus(t, code, http.StatusOK, body)
	if obj(arr(body["items"])[0])["quantity"] != float64(3) {
		t.Fatalf("quantity not updated: %v", body)
	}
	if code, _ = e.do("PUT", path, map[string]any{"quantity": 0}, true); code != http.StatusBadRequest {
		t.Fatalf("qty 0 must be 400, got %d", code)
	}
	if code, _ = e.do("PUT", "/api/v1/cart/items/"+uuid.NewString(), map[string]any{"quantity": 2}, true); code != http.StatusNotFound {
		t.Fatalf("unknown item update: %d", code)
	}

	code, body = e.do("DELETE", path, nil, true)
	requireStatus(t, code, http.StatusOK, body)
	if len(arr(body["items"])) != 0 {
		t.Fatalf("item not removed: %v", body)
	}
	if code, _ = e.do("DELETE", path, nil, true); code != http.StatusNotFound {
		t.Fatalf("second delete: %d", code)
	}
}

func TestCart_Clear(t *testing.T) {
	e := newCartEnv()
	e.do("POST", "/api/v1/cart/items", map[string]any{"product_id": e.product, "quantity": 2}, true)
	code, body := e.do("DELETE", "/api/v1/cart", nil, true)
	requireStatus(t, code, http.StatusOK, body)
	requireKeys(t, body, "message")
	_, body = e.do("GET", "/api/v1/cart", nil, true)
	if len(arr(body["items"])) != 0 || body["item_count"] != float64(0) {
		t.Fatalf("cart not empty: %v", body)
	}
}
