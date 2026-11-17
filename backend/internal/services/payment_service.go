package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"strconv"

	"github.com/google/uuid"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
)

// Errors surfaced by the payment flow.
var (
	ErrPaymentsDisabled = errors.New("payments are not configured")
	ErrInvalidSignature = errors.New("invalid webhook signature")
	ErrOrderNotPayable  = errors.New("order is not awaiting payment")
	ErrMissingIdemKey   = errors.New("idempotency key is required")
)

// JobQueue publishes background jobs.
type JobQueue interface {
	Enqueue(ctx context.Context, id, jobType, payload string) error
}

// PaymentInitiation is the response of Initiate (spec: PaymentInitiation).
type PaymentInitiation struct { //nolint:govet // field order follows the API spec
	OrderID         uuid.UUID `json:"order_id"`
	RazorpayOrderID string    `json:"razorpay_order_id"`
	KeyID           string    `json:"key_id"`
	Amount          int64     `json:"amount"`
	Currency        string    `json:"currency"`
}

// PaymentService runs the Razorpay payment lifecycle.
type PaymentService struct {
	gateway       PaymentGateway // nil when Razorpay is not configured
	orders        repositories.OrderRepository
	payments      repositories.PaymentRepository
	idem          IdempotencyStore
	queue         JobQueue
	webhookSecret string
}

// NewPaymentService returns a PaymentService. gateway may be nil (initiate then answers ErrPaymentsDisabled).
func NewPaymentService(gateway PaymentGateway, orders repositories.OrderRepository, payments repositories.PaymentRepository, idem IdempotencyStore, queue JobQueue, webhookSecret string) *PaymentService {
	return &PaymentService{gateway: gateway, orders: orders, payments: payments, idem: idem, queue: queue, webhookSecret: webhookSecret}
}

// Initiate creates a Razorpay order for one of the user's pending orders. Replaying the
// same idempotency key returns the first response.
func (s *PaymentService) Initiate(ctx context.Context, userID, orderID uuid.UUID, idemKey string) (*PaymentInitiation, error) {
	if s.gateway == nil {
		return nil, ErrPaymentsDisabled
	}
	if len(idemKey) < 8 || len(idemKey) > 128 {
		return nil, ErrMissingIdemKey
	}
	order, err := s.orders.GetByID(ctx, orderID)
	if err != nil {
		return nil, mapNotFound(err)
	}
	if order.UserID != userID {
		return nil, ErrNotFound // do not reveal other users' orders
	}

	key := fmt.Sprintf("idem:payments:initiate:%s:%s:%s", userID, orderID, idemKey)
	cached, done, err := s.idem.Begin(ctx, key)
	if err != nil {
		return nil, err
	}
	if done {
		var out PaymentInitiation
		if err = json.Unmarshal(cached, &out); err != nil {
			return nil, err
		}
		return &out, nil
	}
	out, err := s.create(ctx, order)
	if err != nil {
		s.idem.Abort(ctx, key)
		return nil, err
	}
	if b, merr := json.Marshal(out); merr == nil {
		if cerr := s.idem.Complete(ctx, key, b); cerr != nil {
			slog.Error("payments: store idempotent response", "err", cerr)
		}
	}
	return out, nil
}

func (s *PaymentService) create(ctx context.Context, order *dbsqlc.Order) (*PaymentInitiation, error) {
	if order.Status != dbsqlc.OrderStatusPendingPayment {
		return nil, ErrOrderNotPayable
	}
	total, err := strconv.ParseFloat(order.TotalInr, 64)
	if err != nil {
		return nil, err
	}
	paise := int64(math.Round(total * 100))
	rzpID, err := s.gateway.CreateOrder(ctx, paise, order.OrderNumber)
	if err != nil {
		return nil, err
	}
	if _, err := s.payments.Create(ctx, order.ID, rzpID, order.TotalInr); err != nil {
		return nil, err
	}
	return &PaymentInitiation{OrderID: order.ID, RazorpayOrderID: rzpID, KeyID: s.gateway.KeyID(), Amount: paise, Currency: "INR"}, nil
}

// VerifySignature checks the hex HMAC-SHA256 of body against the webhook secret.
func (s *PaymentService) VerifySignature(body []byte, signature string) bool {
	if s.webhookSecret == "" || signature == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(s.webhookSecret))
	mac.Write(body)
	want := mac.Sum(nil)
	got, err := hex.DecodeString(signature)
	return err == nil && hmac.Equal(got, want)
}

type webhookEvent struct { //nolint:govet // mirrors the Razorpay payload
	Event   string `json:"event"`
	Payload struct {
		Payment struct {
			Entity struct {
				ID      string `json:"id"`
				OrderID string `json:"order_id"`
			} `json:"entity"`
		} `json:"payment"`
	} `json:"payload"`
}

// HandleWebhook verifies and applies a Razorpay event. Unknown events and unknown
// payments are acknowledged (nil) so Razorpay does not retry them forever.
func (s *PaymentService) HandleWebhook(ctx context.Context, body []byte, signature string) error {
	if !s.VerifySignature(body, signature) {
		return ErrInvalidSignature
	}
	var ev webhookEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		return ErrInvalidSignature // a signed but unparseable body is still unusable
	}
	var outcome string
	switch ev.Event {
	case "payment.captured":
		outcome = repositories.PaymentCaptured
	case "payment.failed":
		outcome = repositories.PaymentFailed
	default:
		return nil
	}
	e := ev.Payload.Payment.Entity
	res, err := s.payments.Apply(ctx, e.OrderID, e.ID, outcome, body)
	if errors.Is(err, repositories.ErrPaymentNotFound) {
		slog.Warn("payments: webhook for unknown payment", "provider_order_id", e.OrderID)
		return nil
	}
	if err != nil {
		return err
	}
	// Enqueue even on replay: the job id is idempotent and a previous attempt may have
	// committed the payment but failed to publish.
	if outcome == repositories.PaymentCaptured {
		id := "order_confirmation:" + res.OrderID.String()
		if err := s.queue.Enqueue(ctx, id, "order_confirmation_email", res.OrderID.String()); err != nil {
			return fmt.Errorf("payments: enqueue confirmation: %w", err)
		}
	}
	return nil
}
