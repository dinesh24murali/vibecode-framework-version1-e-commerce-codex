package contract

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
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
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

const webhookSecret = "whsec_test"

type fakeGateway struct{ calls int }

func (g *fakeGateway) CreateOrder(_ context.Context, amountPaise int64, _ string) (string, error) {
	g.calls++
	_ = amountPaise
	return "order_RZP" + uuid.NewString()[:6], nil
}
func (g *fakeGateway) KeyID() string { return "rzp_test_key" }

type fakeIdem struct{ store map[string][]byte }

func (f *fakeIdem) Begin(_ context.Context, key string) ([]byte, bool, error) {
	if v, ok := f.store[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}
func (f *fakeIdem) Complete(_ context.Context, key string, v []byte) error {
	f.store[key] = v
	return nil
}
func (f *fakeIdem) Abort(context.Context, string) {}

type fakeQueue struct{ jobs []string }

func (q *fakeQueue) Enqueue(_ context.Context, id, jobType, payload string) error {
	q.jobs = append(q.jobs, id+"|"+jobType+"|"+payload)
	return nil
}

// fakePayments mirrors the repository semantics (idempotent Apply) on top of fakeOrders.
type fakePayments struct {
	orders   *fakeOrders
	byRzp    map[string]*dbsqlc.Payment
	deducted int
}

func (f *fakePayments) Create(_ context.Context, orderID uuid.UUID, rzp, amount string) (*dbsqlc.Payment, error) {
	p := &dbsqlc.Payment{ID: uuid.New(), OrderID: orderID, Provider: "razorpay", ProviderOrderID: rzp, Status: repositories.PaymentCreated, AmountInr: amount}
	f.byRzp[rzp] = p
	return p, nil
}

func (f *fakePayments) Apply(_ context.Context, rzp, payID, outcome string, _ []byte) (*repositories.PaymentResult, error) {
	p, ok := f.byRzp[rzp]
	if !ok {
		return nil, repositories.ErrPaymentNotFound
	}
	res := &repositories.PaymentResult{OrderID: p.OrderID}
	if p.Status == repositories.PaymentCaptured || p.Status == outcome {
		return res, nil
	}
	p.Status = outcome
	p.ProviderPaymentID = sql.NullString{String: payID, Valid: true}
	st := dbsqlc.OrderStatusPaymentFailed
	if outcome == repositories.PaymentCaptured {
		st = dbsqlc.OrderStatusPaid
		f.deducted++
	}
	_, _ = f.orders.UpdateStatus(context.Background(), p.OrderID, st)
	res.Applied = true
	return res, nil
}

type payEnv struct {
	r      *gin.Engine
	orders *fakeOrders
	pays   *fakePayments
	gw     *fakeGateway
	queue  *fakeQueue
	me     uuid.UUID
	order  uuid.UUID
}

func newPayEnv(configured bool) *payEnv {
	gin.SetMode(gin.TestMode)
	e := &payEnv{me: uuid.New(), orders: &fakeOrders{}, gw: &fakeGateway{}, queue: &fakeQueue{}}
	e.pays = &fakePayments{orders: e.orders, byRzp: map[string]*dbsqlc.Payment{}}
	e.order = uuid.New()
	e.orders.rows = append(e.orders.rows, dbsqlc.Order{
		ID: e.order, UserID: e.me, OrderNumber: "ORD-1", Status: dbsqlc.OrderStatusPendingPayment,
		SubtotalInr: "750.00", TaxInr: "135.00", TotalInr: "885.00", CreatedAt: time.Now(), UpdatedAt: time.Now(),
	})
	var gw services.PaymentGateway
	if configured {
		gw = e.gw
	}
	svc := services.NewPaymentService(gw, e.orders, e.pays, &fakeIdem{store: map[string][]byte{}}, e.queue, webhookSecret)
	h := handlers.NewPaymentHandler(svc)
	r := gin.New()
	r.POST("/api/v1/payments/webhooks/razorpay", h.Webhook)
	g := r.Group("/api/v1", func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "unauthorized"})
			return
		}
		c.Set(auth.ContextUserID, e.me)
	})
	g.POST("/payments/initiate", h.Initiate)
	e.r = r
	return e
}

func (e *payEnv) initiate(orderID uuid.UUID, idemKey string, authed bool) (int, map[string]any) {
	b, _ := json.Marshal(map[string]any{"order_id": orderID})
	req := httptest.NewRequest("POST", "/api/v1/payments/initiate", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	if authed {
		req.Header.Set("Authorization", "Bearer x")
	}
	if idemKey != "" {
		req.Header.Set("Idempotency-Key", idemKey)
	}
	w := httptest.NewRecorder()
	e.r.ServeHTTP(w, req)
	var out map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &out)
	return w.Code, out
}

func sign(body []byte, secret string) string {
	m := hmac.New(sha256.New, []byte(secret))
	m.Write(body)
	return hex.EncodeToString(m.Sum(nil))
}

func (e *payEnv) webhook(event, rzpOrderID, signature string) int {
	body := []byte(`{"event":"` + event + `","payload":{"payment":{"entity":{"id":"pay_1","order_id":"` + rzpOrderID + `"}}}}`)
	if signature == "" {
		signature = sign(body, webhookSecret)
	}
	req := httptest.NewRequest("POST", "/api/v1/payments/webhooks/razorpay", bytes.NewReader(body))
	req.Header.Set("X-Razorpay-Signature", signature)
	w := httptest.NewRecorder()
	e.r.ServeHTTP(w, req)
	return w.Code
}

func (e *payEnv) status() dbsqlc.OrderStatus { return e.orders.rows[0].Status }

func TestPayments_InitiateRequiresAuthAndKey(t *testing.T) {
	e := newPayEnv(true)
	if code, _ := e.initiate(e.order, "idem-key-1", false); code != http.StatusUnauthorized {
		t.Fatalf("unauth: %d", code)
	}
	if code, _ := e.initiate(e.order, "", true); code != http.StatusBadRequest {
		t.Fatalf("missing key: %d", code)
	}
	if code, _ := e.initiate(uuid.New(), "idem-key-1", true); code != http.StatusNotFound {
		t.Fatalf("unknown order: %d", code)
	}
}

func TestPayments_InitiateReturnsRazorpayOrderAndIsIdempotent(t *testing.T) {
	e := newPayEnv(true)
	code, body := e.initiate(e.order, "idem-key-1", true)
	requireStatus(t, code, http.StatusOK, body)
	requireKeys(t, body, "order_id", "razorpay_order_id", "key_id", "amount", "currency")
	if body["key_id"] != "rzp_test_key" || body["amount"] != float64(88500) || body["currency"] != "INR" {
		t.Fatalf("body: %v", body)
	}
	_, again := e.initiate(e.order, "idem-key-1", true)
	if again["razorpay_order_id"] != body["razorpay_order_id"] || e.gw.calls != 1 {
		t.Fatalf("replay must return the cached response (gateway calls=%d)", e.gw.calls)
	}
	if code, _ = e.initiate(e.order, "idem-key-2", true); code != http.StatusOK || e.gw.calls != 2 {
		t.Fatalf("a new key creates a new Razorpay order: %d calls=%d", code, e.gw.calls)
	}
}

func TestPayments_InitiateRejectsOtherUsersAndNonPendingOrders(t *testing.T) {
	e := newPayEnv(true)
	e.orders.rows[0].UserID = uuid.New()
	if code, _ := e.initiate(e.order, "idem-key-1", true); code != http.StatusNotFound {
		t.Fatalf("other user's order: %d", code)
	}
	e.orders.rows[0].UserID = e.me
	e.orders.rows[0].Status = dbsqlc.OrderStatusPaid
	if code, _ := e.initiate(e.order, "idem-key-1", true); code != http.StatusConflict {
		t.Fatalf("paid order: %d", code)
	}
}

func TestPayments_InitiateUnconfiguredIs503(t *testing.T) {
	e := newPayEnv(false)
	if code, _ := e.initiate(e.order, "idem-key-1", true); code != http.StatusServiceUnavailable {
		t.Fatalf("got %d", code)
	}
}

func TestPayments_WebhookRejectsBadSignature(t *testing.T) {
	e := newPayEnv(true)
	_, body := e.initiate(e.order, "idem-key-1", true)
	rzp, _ := body["razorpay_order_id"].(string)
	if code := e.webhook("payment.captured", rzp, "deadbeef"); code != http.StatusBadRequest {
		t.Fatalf("bad signature: %d", code)
	}
	if code := e.webhook("payment.captured", rzp, sign([]byte("other"), "wrong-secret")); code != http.StatusBadRequest {
		t.Fatalf("wrong secret: %d", code)
	}
	if e.status() != dbsqlc.OrderStatusPendingPayment || len(e.queue.jobs) != 0 {
		t.Fatal("nothing may change on a rejected webhook")
	}
}

func TestPayments_WebhookCapturedMarksPaidOnceAndEnqueuesJob(t *testing.T) {
	e := newPayEnv(true)
	_, body := e.initiate(e.order, "idem-key-1", true)
	rzp, _ := body["razorpay_order_id"].(string)
	if code := e.webhook("payment.captured", rzp, ""); code != http.StatusOK {
		t.Fatalf("captured: %d", code)
	}
	if e.status() != dbsqlc.OrderStatusPaid || e.pays.byRzp[rzp].Status != repositories.PaymentCaptured {
		t.Fatalf("order=%s payment=%s", e.status(), e.pays.byRzp[rzp].Status)
	}
	if len(e.queue.jobs) != 1 || e.queue.jobs[0] != "order_confirmation:"+e.order.String()+"|order_confirmation_email|"+e.order.String() {
		t.Fatalf("jobs: %v", e.queue.jobs)
	}
	// Replay: still 200, stock deducted only once, same (idempotent) job id.
	if code := e.webhook("payment.captured", rzp, ""); code != http.StatusOK || e.pays.deducted != 1 {
		t.Fatalf("replay: %d deducted=%d", code, e.pays.deducted)
	}
	// A late failure cannot undo a capture.
	if code := e.webhook("payment.failed", rzp, ""); code != http.StatusOK || e.status() != dbsqlc.OrderStatusPaid {
		t.Fatalf("late failure: %d status=%s", code, e.status())
	}
}

func TestPayments_WebhookFailedMarksPaymentFailed(t *testing.T) {
	e := newPayEnv(true)
	_, body := e.initiate(e.order, "idem-key-1", true)
	rzp, _ := body["razorpay_order_id"].(string)
	if code := e.webhook("payment.failed", rzp, ""); code != http.StatusOK {
		t.Fatalf("failed: %d", code)
	}
	if e.status() != dbsqlc.OrderStatusPaymentFailed || len(e.queue.jobs) != 0 {
		t.Fatalf("status=%s jobs=%v", e.status(), e.queue.jobs)
	}
}

func TestPayments_WebhookIgnoresUnknownEventsAndPayments(t *testing.T) {
	e := newPayEnv(true)
	if code := e.webhook("refund.created", "order_x", ""); code != http.StatusOK {
		t.Fatalf("unknown event: %d", code)
	}
	if code := e.webhook("payment.captured", "order_unknown", ""); code != http.StatusOK {
		t.Fatalf("unknown payment: %d", code)
	}
}
