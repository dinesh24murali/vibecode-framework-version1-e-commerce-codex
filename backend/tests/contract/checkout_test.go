package contract

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/auth"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/handlers"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

type fakeCoupons struct {
	byCode map[string]*dbsqlc.Coupon
	used   map[[2]uuid.UUID]bool
}

func (f *fakeCoupons) GetByCode(_ context.Context, code string) (*dbsqlc.Coupon, error) {
	if c, ok := f.byCode[code]; ok {
		return c, nil
	}
	return nil, sql.ErrNoRows
}

func (f *fakeCoupons) HasUsed(_ context.Context, u, c uuid.UUID) (bool, error) {
	return f.used[[2]uuid.UUID{u, c}], nil
}

// fakeCheckout mimics the transactional repository: stock check, order, coupon use, cart clear.
type fakeCheckout struct { //nolint:govet // field order is for readability
	stock   int32
	coupons *fakeCoupons
	cart    *fakeCart
	placed  []*repositories.OrderDraft
}

func (f *fakeCheckout) Place(_ context.Context, d *repositories.OrderDraft) (*repositories.PlacedOrder, error) {
	for _, it := range d.Items {
		if it.Quantity > f.stock {
			return nil, repositories.ErrInsufficientStock
		}
	}
	for _, it := range d.Items {
		f.stock -= it.Quantity
	}
	if d.CouponID.Valid {
		f.coupons.used[[2]uuid.UUID{d.UserID, d.CouponID.UUID}] = true
	}
	_ = f.cart.Clear(context.Background(), d.UserID)
	f.placed = append(f.placed, d)
	o := &dbsqlc.Order{
		ID: uuid.New(), UserID: d.UserID, OrderNumber: d.OrderNumber, Status: dbsqlc.OrderStatusPendingPayment,
		SubtotalInr: "x", TaxInr: "x", TotalInr: "x", CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	o.SubtotalInr, o.TaxInr, o.TotalInr = fmtMoney(d.Subtotal), fmtMoney(d.GST), fmtMoney(d.Total)
	o.CouponDiscountInr = fmtMoney(d.Discount)
	if d.CouponCode != "" {
		o.CouponCode = sql.NullString{String: d.CouponCode, Valid: true}
	}
	items := make([]dbsqlc.ListOrderItemsWithProductRow, 0, len(d.Items))
	for _, it := range d.Items {
		items = append(items, dbsqlc.ListOrderItemsWithProductRow{
			ID: uuid.New(), ProductID: uuid.NullUUID{UUID: it.ProductID, Valid: true}, ProductTitle: it.Title,
			Quantity: it.Quantity, UnitPriceInr: fmtMoney(it.UnitPrice), LineTotalInr: fmtMoney(it.LineTotal), ProductSlug: it.SKU,
		})
	}
	return &repositories.PlacedOrder{Order: o, Items: items, Address: &dbsqlc.OrderAddress{ID: uuid.New(), OrderID: o.ID, Line1: d.Address.Line1, City: d.Address.City}}, nil
}

func fmtMoney(f float64) string {
	b, _ := json.Marshal(math.Round(f*100) / 100)
	return string(b)
}

type checkoutEnv struct { //nolint:govet // field order is for readability
	r       *gin.Engine
	product uuid.UUID
	address uuid.UUID
	repo    *fakeCheckout
	cart    *fakeCart
}

func newCheckoutEnv(stock int32) *checkoutEnv {
	gin.SetMode(gin.TestMode)
	now := time.Now()
	me := uuid.New()
	cat := dbsqlc.Category{ID: uuid.New(), Slug: "journals", Name: "Journals", CreatedAt: now, UpdatedAt: now}
	fp := &fakeProducts{product: dbsqlc.GetActiveProductByIDRow{
		ID: uuid.New(), CategoryID: cat.ID, Slug: "mandala-journal", Title: "Mandala Journal", PriceInr: "500.00",
		DiscountPercent: sql.NullInt32{Int32: 10, Valid: true}, IsActive: true, CreatedAt: now, UpdatedAt: now,
	}}
	attrSvc := services.NewAttributeService(&fakeAttrs{values: map[uuid.UUID][]dbsqlc.ProductAttributeValue{}})
	products := services.NewProductService(fp, &fakeCategories{rows: []dbsqlc.Category{cat}}, attrSvc)
	cart := &fakeCart{}
	cartSvc := services.NewCartService(cart, products)

	superhit := &dbsqlc.Coupon{ID: uuid.New(), Code: "SUPERHIT", Type: "flat", Flat: sql.NullInt32{Int32: 150, Valid: true}, IsEnable: true}
	firsttime := &dbsqlc.Coupon{ID: uuid.New(), Code: "FIRSTTIME", Type: "percentage", Percent: sql.NullInt32{Int32: 10, Valid: true}, IsEnable: true}
	expired := &dbsqlc.Coupon{ID: uuid.New(), Code: "OLD", Type: "flat", Flat: sql.NullInt32{Int32: 10, Valid: true}, IsEnable: true,
		CouponExpiry: sql.NullTime{Time: now.Add(-time.Hour), Valid: true}}
	coupons := &fakeCoupons{byCode: map[string]*dbsqlc.Coupon{"SUPERHIT": superhit, "FIRSTTIME": firsttime, "OLD": expired}, used: map[[2]uuid.UUID]bool{}}

	addrs := &fakeAddresses{}
	a, _ := addrs.Create(context.Background(), dbsqlc.CreateAddressParams{UserID: me, Line1: "42 Rose St", City: "Chennai", State: "TN", PostalCode: "600001"})
	repo := &fakeCheckout{stock: stock, coupons: coupons, cart: cart}
	h := handlers.NewCheckoutHandler(services.NewCheckoutService(cartSvc, services.NewCouponService(coupons), addrs, repo))

	r := gin.New()
	g := r.Group("/api/v1", func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "unauthorized"})
			return
		}
		c.Set(auth.ContextUserID, me)
	})
	g.POST("/checkout/calculate", h.Calculate)
	g.POST("/checkout/coupons/validate", h.ValidateCoupon)
	g.POST("/checkout/orders", h.PlaceOrder)
	e := &checkoutEnv{r: r, product: fp.product.ID, address: a.ID, repo: repo, cart: cart}
	e.addToCart(me, 2) // 2 x 450 = 900
	return e
}

func (e *checkoutEnv) addToCart(user uuid.UUID, qty int32) {
	_ = e.cart.Upsert(context.Background(), user, e.product, qty)
}

func (e *checkoutEnv) do(path string, body any, authed bool) (int, map[string]any) {
	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(body)
	req := httptest.NewRequest("POST", path, &buf)
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

func near(t *testing.T, name string, got any, want float64) {
	t.Helper()
	f, _ := got.(float64)
	if math.Abs(f-want) > 0.005 {
		t.Fatalf("%s = %v, want %v", name, got, want)
	}
}

func TestCheckout_RequiresAuth(t *testing.T) {
	e := newCheckoutEnv(10)
	for _, p := range []string{"/api/v1/checkout/calculate", "/api/v1/checkout/coupons/validate", "/api/v1/checkout/orders"} {
		if code, _ := e.do(p, map[string]any{}, false); code != http.StatusUnauthorized {
			t.Fatalf("%s: %d", p, code)
		}
	}
}

func TestCheckout_CalculateNoCoupon(t *testing.T) {
	e := newCheckoutEnv(10)
	code, body := e.do("/api/v1/checkout/calculate", map[string]any{"address_id": e.address}, true)
	requireStatus(t, code, http.StatusOK, body)
	requireKeys(t, body, "items", "subtotal", "coupon_discount", "taxable_amount", "gst_amount", "total", "address")
	near(t, "subtotal", body["subtotal"], 900)
	near(t, "coupon_discount", body["coupon_discount"], 0)
	near(t, "gst", body["gst_amount"], 162) // 18% of 900
	near(t, "total", body["total"], 1062)
	if body["coupon_code"] != nil {
		t.Fatalf("coupon_code should be null: %v", body["coupon_code"])
	}
}

func TestCheckout_CalculateFlatCouponAppliedBeforeTax(t *testing.T) {
	e := newCheckoutEnv(10)
	code, body := e.do("/api/v1/checkout/calculate", map[string]any{"address_id": e.address, "coupon_code": "superhit"}, true)
	requireStatus(t, code, http.StatusOK, body)
	near(t, "coupon_discount", body["coupon_discount"], 150)
	near(t, "taxable", body["taxable_amount"], 750)
	near(t, "gst", body["gst_amount"], 135) // tax on the post-coupon amount
	near(t, "total", body["total"], 885)
	if body["coupon_code"] != "SUPERHIT" {
		t.Fatalf("coupon_code: %v", body["coupon_code"])
	}
}

func TestCheckout_CalculatePercentageCoupon(t *testing.T) {
	e := newCheckoutEnv(10)
	code, body := e.do("/api/v1/checkout/calculate", map[string]any{"address_id": e.address, "coupon_code": "FIRSTTIME"}, true)
	requireStatus(t, code, http.StatusOK, body)
	near(t, "coupon_discount", body["coupon_discount"], 90)
	near(t, "taxable", body["taxable_amount"], 810)
	near(t, "gst", body["gst_amount"], 145.8)
	near(t, "total", body["total"], 955.8)
}

func TestCheckout_CalculateEdgeCases(t *testing.T) {
	e := newCheckoutEnv(10)
	if code, _ := e.do("/api/v1/checkout/calculate", map[string]any{}, true); code != http.StatusBadRequest {
		t.Fatalf("missing address: %d", code)
	}
	if code, _ := e.do("/api/v1/checkout/calculate", map[string]any{"address_id": uuid.New()}, true); code != http.StatusUnprocessableEntity {
		t.Fatalf("unknown address: %d", code)
	}
	for _, c := range []string{"NOPE", "OLD"} {
		if code, _ := e.do("/api/v1/checkout/calculate", map[string]any{"address_id": e.address, "coupon_code": c}, true); code != http.StatusUnprocessableEntity {
			t.Fatalf("coupon %s: %d", c, code)
		}
	}
}

func TestCheckout_ValidateCoupon(t *testing.T) {
	e := newCheckoutEnv(10)
	code, body := e.do("/api/v1/checkout/coupons/validate", map[string]any{"coupon_code": "SUPERHIT"}, true)
	requireStatus(t, code, http.StatusOK, body)
	requireKeys(t, body, "is_valid", "coupon_code", "discount_type", "discount_value", "message")
	if body["is_valid"] != true || body["discount_type"] != "flat" {
		t.Fatalf("flat: %v", body)
	}
	near(t, "value", body["discount_value"], 150)

	_, body = e.do("/api/v1/checkout/coupons/validate", map[string]any{"coupon_code": "FIRSTTIME"}, true)
	if body["discount_type"] != "percentage" {
		t.Fatalf("percentage: %v", body)
	}
	near(t, "value", body["discount_value"], 10)

	// Unknown and expired codes are a normal 200 with is_valid=false (spec has no 404 here).
	for _, c := range []string{"NOPE", "OLD"} {
		code, body = e.do("/api/v1/checkout/coupons/validate", map[string]any{"coupon_code": c}, true)
		requireStatus(t, code, http.StatusOK, body)
		if body["is_valid"] != false || body["message"] == "" {
			t.Fatalf("%s: %v", c, body)
		}
	}
	if code, _ = e.do("/api/v1/checkout/coupons/validate", map[string]any{}, true); code != http.StatusBadRequest {
		t.Fatalf("missing code: %d", code)
	}
}

func TestCheckout_PlaceOrderDecrementsStockAndClearsCart(t *testing.T) {
	e := newCheckoutEnv(5)
	code, body := e.do("/api/v1/checkout/orders", map[string]any{"address_id": e.address, "coupon_code": "SUPERHIT"}, true)
	requireStatus(t, code, http.StatusCreated, body)
	requireKeys(t, body, "id", "order_number", "status", "items", "shipping_address", "subtotal", "coupon_code", "coupon_discount", "taxable_amount", "gst_amount", "total")
	near(t, "total", body["total"], 885)
	near(t, "taxable", body["taxable_amount"], 750)
	if body["coupon_code"] != "SUPERHIT" || len(arr(body["items"])) != 1 {
		t.Fatalf("order: %v", body)
	}
	if e.repo.stock != 3 {
		t.Fatalf("stock = %d, want 3", e.repo.stock)
	}
	if len(e.cart.rows) != 0 {
		t.Fatal("cart should be cleared")
	}
	// Cart is now empty.
	if code, _ = e.do("/api/v1/checkout/orders", map[string]any{"address_id": e.address}, true); code != http.StatusBadRequest {
		t.Fatalf("empty cart: %d", code)
	}
}

func TestCheckout_PlaceOrderRejectsReusedCouponAndLowStock(t *testing.T) {
	e := newCheckoutEnv(5)
	if code, _ := e.do("/api/v1/checkout/orders", map[string]any{"address_id": e.address, "coupon_code": "FIRSTTIME"}, true); code != http.StatusCreated {
		t.Fatal("first use should succeed")
	}
	// Refill the cart for the same user and try the coupon again.
	for _, r := range e.repo.placed {
		e.addToCart(r.UserID, 1)
	}
	if code, _ := e.do("/api/v1/checkout/orders", map[string]any{"address_id": e.address, "coupon_code": "FIRSTTIME"}, true); code != http.StatusUnprocessableEntity {
		t.Fatalf("reused coupon: %d", code)
	}

	low := newCheckoutEnv(1) // cart holds 2
	code, body := low.do("/api/v1/checkout/orders", map[string]any{"address_id": low.address}, true)
	if code != http.StatusUnprocessableEntity {
		t.Fatalf("insufficient stock: %d %v", code, body)
	}
	if low.repo.stock != 1 {
		t.Fatalf("stock must be untouched, got %d", low.repo.stock)
	}
}
