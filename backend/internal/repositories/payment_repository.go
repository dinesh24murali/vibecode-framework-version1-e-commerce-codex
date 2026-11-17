package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/sqlc-dev/pqtype"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
)

// Payment statuses stored in payments.status.
const (
	PaymentCreated  = "created"
	PaymentCaptured = "captured"
	PaymentFailed   = "failed"
)

// ErrPaymentNotFound is returned when no payment matches the gateway order id.
var ErrPaymentNotFound = errors.New("payment not found")

// PaymentResult reports what Apply did.
type PaymentResult struct {
	OrderID uuid.UUID
	// Applied is false when the event was a replay or arrived out of order and nothing changed.
	Applied bool
}

// PaymentRepository persists payments and applies gateway outcomes.
type PaymentRepository interface {
	Create(ctx context.Context, orderID uuid.UUID, providerOrderID, amountInr string) (*dbsqlc.Payment, error)
	// Apply records a captured or failed outcome idempotently. Captured moves the order
	// to paid and deducts stock (on_hand and reserved); failed moves it to payment_failed
	// and leaves the reservation for the stock_release job.
	Apply(ctx context.Context, providerOrderID, providerPaymentID, outcome string, payload []byte) (*PaymentResult, error)
}

type pgPaymentRepository struct {
	db *sql.DB
	q  *dbsqlc.Queries
}

// NewPaymentRepository returns a Postgres-backed PaymentRepository.
func NewPaymentRepository(pool *pgxpool.Pool) PaymentRepository {
	db := stdlib.OpenDBFromPool(pool)
	return &pgPaymentRepository{db: db, q: dbsqlc.New(db)}
}

func (r *pgPaymentRepository) Create(ctx context.Context, orderID uuid.UUID, providerOrderID, amountInr string) (*dbsqlc.Payment, error) {
	p, err := r.q.CreatePayment(ctx, dbsqlc.CreatePaymentParams{
		OrderID: orderID, Provider: "razorpay", ProviderOrderID: providerOrderID, Status: PaymentCreated, AmountInr: amountInr,
	})
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *pgPaymentRepository) Apply(ctx context.Context, providerOrderID, providerPaymentID, outcome string, payload []byte) (*PaymentResult, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	q := r.q.WithTx(tx)

	p, err := q.GetPaymentForUpdate(ctx, providerOrderID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPaymentNotFound
	}
	if err != nil {
		return nil, err
	}
	res := &PaymentResult{OrderID: p.OrderID}
	// A captured payment is final; a repeated failure is a no-op.
	if p.Status == PaymentCaptured || p.Status == outcome {
		return res, nil
	}

	raw := pqtype.NullRawMessage{}
	if json.Valid(payload) {
		raw = pqtype.NullRawMessage{RawMessage: payload, Valid: true}
	}
	if _, err = q.UpdatePaymentStatus(ctx, dbsqlc.UpdatePaymentStatusParams{
		ID: p.ID, Status: outcome, ProviderPaymentID: sql.NullString{String: providerPaymentID, Valid: providerPaymentID != ""}, GatewayPayload: raw,
	}); err != nil {
		return nil, err
	}

	if outcome == PaymentCaptured {
		items, ierr := q.ListOrderItems(ctx, p.OrderID)
		if ierr != nil {
			return nil, ierr
		}
		for _, it := range items {
			if !it.ProductID.Valid {
				continue
			}
			if _, err = q.DeductStock(ctx, dbsqlc.DeductStockParams{ProductID: it.ProductID.UUID, OnHand: it.Quantity}); err != nil {
				return nil, err
			}
		}
		_, err = q.UpdateOrderStatus(ctx, dbsqlc.UpdateOrderStatusParams{ID: p.OrderID, Status: dbsqlc.OrderStatusPaid})
	} else {
		_, err = q.UpdateOrderStatus(ctx, dbsqlc.UpdateOrderStatusParams{ID: p.OrderID, Status: dbsqlc.OrderStatusPaymentFailed})
	}
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	res.Applied = true
	return res, nil
}
