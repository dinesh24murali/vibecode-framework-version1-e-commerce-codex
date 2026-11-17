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

// ── in-memory fakes (the real handlers + service run on top of them) ─────────

type fakeUsers struct{ users map[uuid.UUID]*dbsqlc.User }

func (f *fakeUsers) GetByID(_ context.Context, id uuid.UUID) (*dbsqlc.User, error) {
	if u, ok := f.users[id]; ok {
		return u, nil
	}
	return nil, sql.ErrNoRows
}
func (f *fakeUsers) GetByEmail(_ context.Context, email string) (*dbsqlc.User, error) {
	for _, u := range f.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, sql.ErrNoRows
}
func (f *fakeUsers) Create(context.Context, dbsqlc.CreateUserParams) (*dbsqlc.User, error) {
	return nil, nil
}
func (f *fakeUsers) UpdatePasswordHash(context.Context, uuid.UUID, string) error { return nil }
func (f *fakeUsers) UpdateName(_ context.Context, id uuid.UUID, n string) (*dbsqlc.User, error) {
	u := f.users[id]
	u.FullName = n
	return u, nil
}
func (f *fakeUsers) UpdateEmail(_ context.Context, id uuid.UUID, e string) (*dbsqlc.User, error) {
	u := f.users[id]
	u.Email = e
	return u, nil
}

type fakeAddresses struct{ rows []dbsqlc.Address }

func (f *fakeAddresses) ListByUser(_ context.Context, uid uuid.UUID) ([]dbsqlc.Address, error) {
	var out []dbsqlc.Address
	for _, a := range f.rows {
		if a.UserID == uid {
			out = append(out, a)
		}
	}
	return out, nil
}
func (f *fakeAddresses) GetForUser(_ context.Context, id, uid uuid.UUID) (*dbsqlc.Address, error) {
	for i := range f.rows {
		if f.rows[i].ID == id && f.rows[i].UserID == uid {
			return &f.rows[i], nil
		}
	}
	return nil, sql.ErrNoRows
}
func (f *fakeAddresses) Create(_ context.Context, p dbsqlc.CreateAddressParams) (*dbsqlc.Address, error) {
	if p.IsDefault {
		for i := range f.rows {
			f.rows[i].IsDefault = false
		}
	}
	a := dbsqlc.Address{
		ID: uuid.New(), UserID: p.UserID, Line1: p.Line1, Line2: p.Line2, Landmark: p.Landmark,
		City: p.City, State: p.State, PostalCode: p.PostalCode, CompanyName: p.CompanyName,
		GstNumber: p.GstNumber, PhoneE164: p.PhoneE164, IsDefault: p.IsDefault,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	f.rows = append(f.rows, a)
	return &a, nil
}
func (f *fakeAddresses) Update(_ context.Context, p dbsqlc.UpdateAddressParams) (*dbsqlc.Address, error) {
	for i := range f.rows {
		if f.rows[i].ID == p.ID && f.rows[i].UserID == p.UserID {
			f.rows[i].City, f.rows[i].GstNumber = p.City, p.GstNumber
			return &f.rows[i], nil
		}
	}
	return nil, sql.ErrNoRows
}
func (f *fakeAddresses) Delete(_ context.Context, id, uid uuid.UUID) error {
	for i := range f.rows {
		if f.rows[i].ID == id && f.rows[i].UserID == uid {
			f.rows = append(f.rows[:i], f.rows[i+1:]...)
			return nil
		}
	}
	return nil
}

type fakeOrders struct{ rows []dbsqlc.Order }

func (f *fakeOrders) ListByUser(_ context.Context, uid uuid.UUID, st []dbsqlc.OrderStatus, limit, offset int32) ([]dbsqlc.Order, int64, error) {
	var all []dbsqlc.Order
	for _, o := range f.rows {
		if o.UserID != uid {
			continue
		}
		if len(st) > 0 {
			match := false
			for _, s := range st {
				match = match || s == o.Status
			}
			if !match {
				continue
			}
		}
		all = append(all, o)
	}
	total := int64(len(all))
	lo := min(int(offset), len(all))
	hi := min(lo+int(limit), len(all))
	return all[lo:hi], total, nil
}
func (f *fakeOrders) GetByID(_ context.Context, id uuid.UUID) (*dbsqlc.Order, error) {
	for i := range f.rows {
		if f.rows[i].ID == id {
			return &f.rows[i], nil
		}
	}
	return nil, sql.ErrNoRows
}
func (f *fakeOrders) Items(context.Context, uuid.UUID) ([]dbsqlc.ListOrderItemsWithProductRow, error) {
	return []dbsqlc.ListOrderItemsWithProductRow{{
		ID: uuid.New(), ProductTitle: "Book", ProductSlug: "book", Quantity: 2,
		UnitPriceInr: "100.00", LineTotalInr: "200.00",
	}}, nil
}
func (f *fakeOrders) Address(context.Context, uuid.UUID) (*dbsqlc.OrderAddress, error) {
	return nil, nil
}
func (f *fakeOrders) UpdateStatus(_ context.Context, id uuid.UUID, s dbsqlc.OrderStatus) (*dbsqlc.Order, error) {
	for i := range f.rows {
		if f.rows[i].ID == id {
			f.rows[i].Status = s
			return &f.rows[i], nil
		}
	}
	return nil, sql.ErrNoRows
}

type fakeOTP struct{}

func (fakeOTP) GenerateAndStore(context.Context, string, string) (string, error) {
	return "123456", nil
}
func (fakeOTP) Verify(_ context.Context, _, _, otp string) (bool, error) { return otp == "123456", nil }

// ── harness ──────────────────────────────────────────────────────────────────

type acct struct {
	r      *gin.Engine
	orders *fakeOrders
	users  *fakeUsers
	me     uuid.UUID
	other  uuid.UUID
}

func newAcct() *acct {
	gin.SetMode(gin.TestMode)
	a := &acct{me: uuid.New(), other: uuid.New()}
	a.users = &fakeUsers{users: map[uuid.UUID]*dbsqlc.User{
		a.me:    {ID: a.me, Email: "me@example.com", FullName: "Ann Lee", Role: dbsqlc.UserRoleCustomer},
		a.other: {ID: a.other, Email: "taken@example.com", FullName: "Bob", Role: dbsqlc.UserRoleCustomer},
	}}
	a.orders = &fakeOrders{}
	svc := services.NewAccountService(a.users, &fakeAddresses{}, a.orders, fakeOTP{})
	h := handlers.NewAccountHandler(svc)

	r := gin.New()
	g := r.Group("/api/v1", func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "unauthorized"})
			return
		}
		c.Set(auth.ContextUserID, a.me)
	})
	g.GET("/account/profile", h.GetProfile)
	g.PUT("/account/profile", h.UpdateProfile)
	g.POST("/account/profile/email-change/request", h.RequestEmailChange)
	g.POST("/account/profile/email-change/confirm", h.ConfirmEmailChange)
	g.GET("/account/addresses", h.ListAddresses)
	g.POST("/account/addresses", h.CreateAddress)
	g.PUT("/account/addresses/:id", h.UpdateAddress)
	g.DELETE("/account/addresses/:id", h.DeleteAddress)
	g.GET("/orders", h.ListOrders)
	g.GET("/orders/:id", h.GetOrder)
	g.PUT("/orders/:id/cancel", h.CancelOrder)
	a.r = r
	return a
}

func (a *acct) do(method, path string, body any) (int, map[string]any) {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Authorization", "Bearer x")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	a.r.ServeHTTP(w, req)
	var out map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &out)
	return w.Code, out
}

func (a *acct) addOrder(status dbsqlc.OrderStatus, owner uuid.UUID) uuid.UUID {
	id := uuid.New()
	a.orders.rows = append(a.orders.rows, dbsqlc.Order{
		ID: id, UserID: owner, OrderNumber: "ORD-1", Status: status,
		SubtotalInr: "200.00", TaxInr: "10.00", TotalInr: "210.00", CreatedAt: time.Now(), UpdatedAt: time.Now(),
	})
	return id
}

func requireStatus(t *testing.T, got, want int, body map[string]any) {
	t.Helper()
	if got != want {
		t.Fatalf("status %d, want %d (body %v)", got, want, body)
	}
}

// ── tests ────────────────────────────────────────────────────────────────────

func TestAccount_Unauthorized(t *testing.T) {
	a := newAcct()
	w := httptest.NewRecorder()
	a.r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/account/profile", nil))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("got %d", w.Code)
	}
}

func TestAccount_Profile(t *testing.T) {
	a := newAcct()
	code, body := a.do("GET", "/api/v1/account/profile", nil)
	requireStatus(t, code, 200, body)
	if body["first_name"] != "Ann" || body["last_name"] != "Lee" || body["email"] != "me@example.com" {
		t.Fatalf("profile %v", body)
	}
	code, body = a.do("PUT", "/api/v1/account/profile", gin.H{"first_name": "Anne", "last_name": "Mo Lee"})
	requireStatus(t, code, 200, body)
	if body["first_name"] != "Anne" || body["last_name"] != "Mo Lee" {
		t.Fatalf("updated %v", body)
	}
	code, body = a.do("PUT", "/api/v1/account/profile", gin.H{"first_name": "x"})
	requireStatus(t, code, 400, body)
	if body["message"] == nil {
		t.Fatal("error body must follow ErrorResponse {code,message}")
	}
}

func TestAccount_EmailChange(t *testing.T) {
	a := newAcct()
	code, body := a.do("POST", "/api/v1/account/profile/email-change/request", gin.H{"new_email": "new@example.com"})
	requireStatus(t, code, 200, body)
	code, body = a.do("POST", "/api/v1/account/profile/email-change/confirm", gin.H{"new_email": "new@example.com", "otp": "000000"})
	requireStatus(t, code, 422, body)
	if a.users.users[a.me].Email != "me@example.com" {
		t.Fatal("email must not change before OTP verification")
	}
	code, body = a.do("POST", "/api/v1/account/profile/email-change/confirm", gin.H{"new_email": "taken@example.com", "otp": "123456"})
	requireStatus(t, code, 409, body)
	code, body = a.do("POST", "/api/v1/account/profile/email-change/confirm", gin.H{"new_email": "new@example.com", "otp": "123456"})
	requireStatus(t, code, 200, body)
	if body["email"] != "new@example.com" {
		t.Fatalf("email %v", body["email"])
	}
}

func TestAccount_Addresses(t *testing.T) {
	a := newAcct()
	in := gin.H{
		"address_line1": "42 Rose St", "address_line2": "Flat 3B", "landmark": "Mall", "city": "Chennai",
		"state": "TN", "pincode": "600001", "company_name": "Acme", "gst_number": "33AAAAA0000A1Z5",
		"phone": "+919876543210", "is_default": true,
	}
	code, body := a.do("POST", "/api/v1/account/addresses", in)
	requireStatus(t, code, 201, body)
	for _, k := range []string{"id", "address_line1", "address_line2", "landmark", "city", "state", "pincode",
		"company_name", "gst_number", "phone", "is_default", "created_at", "updated_at"} {
		if _, ok := body[k]; !ok {
			t.Errorf("missing field %s", k)
		}
	}
	if body["gst_number"] != "33AAAAA0000A1Z5" || body["landmark"] != "Mall" {
		t.Fatalf("all fields must persist: %v", body)
	}
	id, _ := body["id"].(string)

	in["city"] = "Madurai"
	code, body = a.do("PUT", "/api/v1/account/addresses/"+id, in)
	requireStatus(t, code, 200, body)
	if body["city"] != "Madurai" {
		t.Fatalf("update %v", body)
	}
	code, body = a.do("GET", "/api/v1/account/addresses", nil)
	requireStatus(t, code, 200, body)
	if len(arr(body["data"])) != 1 {
		t.Fatalf("list %v", body)
	}
	code, body = a.do("POST", "/api/v1/account/addresses", gin.H{"city": "x"})
	requireStatus(t, code, 400, body)
	code, body = a.do("PUT", "/api/v1/account/addresses/"+uuid.NewString(), in)
	requireStatus(t, code, 404, body)
	code, body = a.do("DELETE", "/api/v1/account/addresses/"+id, nil)
	requireStatus(t, code, 200, body)
	code, body = a.do("DELETE", "/api/v1/account/addresses/"+id, nil)
	requireStatus(t, code, 404, body)
}

func TestOrders_ListGetCancel(t *testing.T) {
	a := newAcct()
	open := a.addOrder(dbsqlc.OrderStatusPaid, a.me)
	shipped := a.addOrder(dbsqlc.OrderStatusShipped, a.me)
	delivered := a.addOrder(dbsqlc.OrderStatusDelivered, a.me)
	foreign := a.addOrder(dbsqlc.OrderStatusPaid, a.other)

	code, body := a.do("GET", "/api/v1/orders?limit=2", nil)
	requireStatus(t, code, 200, body)
	meta := obj(body["meta"])
	if meta["total"] != float64(3) || meta["total_pages"] != float64(2) || meta["limit"] != float64(2) {
		t.Fatalf("meta %v", meta)
	}
	code, body = a.do("GET", "/api/v1/orders?status=fulfilled", nil)
	requireStatus(t, code, 200, body)
	rows := arr(body["data"])
	if len(rows) != 1 || obj(rows[0])["status"] != "fulfilled" {
		t.Fatalf("filter %v", rows)
	}
	code, body = a.do("GET", "/api/v1/orders?status=bogus", nil)
	requireStatus(t, code, 400, body)

	code, body = a.do("GET", "/api/v1/orders/"+open.String(), nil)
	requireStatus(t, code, 200, body)
	if body["status"] != "open" || body["total"] != float64(210) || body["gst_amount"] != float64(10) {
		t.Fatalf("order %v", body)
	}
	if items := arr(body["items"]); len(items) != 1 || obj(items[0])["product_slug"] != "book" {
		t.Fatalf("items %v", body["items"])
	}
	code, body = a.do("GET", "/api/v1/orders/"+foreign.String(), nil)
	requireStatus(t, code, 403, body)
	code, body = a.do("GET", "/api/v1/orders/"+uuid.NewString(), nil)
	requireStatus(t, code, 404, body)

	code, body = a.do("PUT", "/api/v1/orders/"+open.String()+"/cancel", nil)
	requireStatus(t, code, 200, body)
	if body["status"] != "cancelled" {
		t.Fatalf("cancel %v", body)
	}
	for _, id := range []uuid.UUID{shipped, delivered, open} {
		code, body = a.do("PUT", "/api/v1/orders/"+id.String()+"/cancel", nil)
		requireStatus(t, code, 400, body)
	}
	code, body = a.do("PUT", "/api/v1/orders/"+foreign.String()+"/cancel", nil)
	requireStatus(t, code, 403, body)
}

func arr(v any) []any {
	a, _ := v.([]any)
	return a
}

func obj(v any) map[string]any {
	m, _ := v.(map[string]any)
	return m
}
