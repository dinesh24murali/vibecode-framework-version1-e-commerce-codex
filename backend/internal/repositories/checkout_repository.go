package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
)

// ErrInsufficientStock is returned when a product cannot cover the requested quantity.
var ErrInsufficientStock = errors.New("insufficient stock")

// ErrCouponAlreadyUsed is returned when the user has already redeemed the coupon.
var ErrCouponAlreadyUsed = errors.New("coupon already used")

// DraftItem is one order line, priced at checkout time.
type DraftItem struct { //nolint:govet // field order is for readability
	ProductID uuid.UUID
	Title     string
	SKU       string
	Quantity  int32
	UnitPrice float64
	LineTotal float64
}

// OrderDraft is a fully priced order ready to be persisted (built by the OrderBuilder).
type OrderDraft struct { //nolint:govet // field order is for readability
	UserID      uuid.UUID
	OrderNumber string
	Items       []DraftItem
	Address     dbsqlc.CreateOrderAddressParams // OrderID is filled in on insert
	CouponID    uuid.NullUUID
	CouponCode  string
	Subtotal    float64
	Discount    float64
	Taxable     float64
	GST         float64
	Total       float64
}

// PlacedOrder is the persisted result of Place.
type PlacedOrder struct { //nolint:govet // field order is for readability
	Order   *dbsqlc.Order
	Items   []dbsqlc.ListOrderItemsWithProductRow
	Address *dbsqlc.OrderAddress
}

// CheckoutRepository persists an order atomically.
type CheckoutRepository interface {
	// Place reserves stock, inserts the order, its items and address, records the
	// coupon redemption and empties the cart, all in one transaction.
	Place(ctx context.Context, d *OrderDraft) (*PlacedOrder, error)
}

type pgCheckoutRepository struct {
	db *sql.DB
	q  *dbsqlc.Queries
}

// NewCheckoutRepository returns a Postgres-backed CheckoutRepository.
func NewCheckoutRepository(pool *pgxpool.Pool) CheckoutRepository {
	db := stdlib.OpenDBFromPool(pool)
	return &pgCheckoutRepository{db: db, q: dbsqlc.New(db)}
}

func money(f float64) string { return fmt.Sprintf("%.2f", f) }

func (r *pgCheckoutRepository) Place(ctx context.Context, d *OrderDraft) (*PlacedOrder, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	q := r.q.WithTx(tx)

	// Reserve stock first; the conditional UPDATE makes the check-and-decrement atomic.
	for _, it := range d.Items {
		n, rerr := q.ReserveStockIfAvailable(ctx, dbsqlc.ReserveStockIfAvailableParams{ProductID: it.ProductID, Reserved: it.Quantity})
		if rerr != nil {
			return nil, rerr
		}
		if n == 0 {
			return nil, fmt.Errorf("%w: %s", ErrInsufficientStock, it.Title)
		}
	}

	var couponCode sql.NullString
	if d.CouponCode != "" {
		couponCode = sql.NullString{String: d.CouponCode, Valid: true}
	}
	order, err := q.CreateOrder(ctx, dbsqlc.CreateOrderParams{
		UserID: d.UserID, OrderNumber: d.OrderNumber, Status: dbsqlc.OrderStatusPendingPayment,
		SubtotalInr: money(d.Subtotal), ShippingInr: "0.00", TaxInr: money(d.GST), TotalInr: money(d.Total),
		CurrencyCode: "INR", CouponCode: couponCode, CouponDiscountInr: money(d.Discount),
	})
	if err != nil {
		return nil, err
	}
	for _, it := range d.Items {
		if _, err = q.CreateOrderItem(ctx, dbsqlc.CreateOrderItemParams{
			OrderID: order.ID, ProductID: uuid.NullUUID{UUID: it.ProductID, Valid: true}, ProductTitle: it.Title,
			ProductSku: it.SKU, Quantity: it.Quantity, UnitPriceInr: money(it.UnitPrice), LineTotalInr: money(it.LineTotal),
		}); err != nil {
			return nil, err
		}
	}
	addr := d.Address
	addr.OrderID = order.ID
	if _, err = q.CreateOrderAddress(ctx, addr); err != nil {
		return nil, err
	}
	if d.CouponID.Valid {
		if _, err = q.CreateCouponUser(ctx, dbsqlc.CreateCouponUserParams{UserID: d.UserID, CouponID: d.CouponID.UUID}); err != nil {
			var pg *pgconn.PgError
			if errors.As(err, &pg) && pg.Code == "23505" {
				return nil, ErrCouponAlreadyUsed
			}
			return nil, err
		}
	}
	if err = q.ClearCart(ctx, d.UserID); err != nil {
		return nil, err
	}
	items, err := q.ListOrderItemsWithProduct(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	oa, err := q.GetOrderAddress(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &PlacedOrder{Order: &order, Items: items, Address: &oa}, nil
}
