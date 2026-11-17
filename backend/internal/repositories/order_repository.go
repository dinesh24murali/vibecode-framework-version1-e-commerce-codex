package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
)

// OrderRepository is the read/cancel persistence interface for a customer's orders.
type OrderRepository interface {
	ListByUser(ctx context.Context, userID uuid.UUID, statuses []dbsqlc.OrderStatus, limit, offset int32) ([]dbsqlc.Order, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dbsqlc.Order, error)
	Items(ctx context.Context, orderID uuid.UUID) ([]dbsqlc.ListOrderItemsWithProductRow, error)
	// Address returns the order's shipping address, or nil when none is stored.
	Address(ctx context.Context, orderID uuid.UUID) (*dbsqlc.OrderAddress, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status dbsqlc.OrderStatus) (*dbsqlc.Order, error)
}

type pgOrderRepository struct {
	q *dbsqlc.Queries
}

// NewOrderRepository returns a Postgres-backed OrderRepository.
func NewOrderRepository(pool *pgxpool.Pool) OrderRepository {
	return &pgOrderRepository{q: dbsqlc.New(stdlib.OpenDBFromPool(pool))}
}

func (r *pgOrderRepository) ListByUser(ctx context.Context, userID uuid.UUID, statuses []dbsqlc.OrderStatus, limit, offset int32) ([]dbsqlc.Order, int64, error) {
	if statuses == nil {
		statuses = []dbsqlc.OrderStatus{}
	}
	total, err := r.q.CountOrdersByUserFiltered(ctx, dbsqlc.CountOrdersByUserFilteredParams{UserID: userID, Column2: statuses})
	if err != nil {
		return nil, 0, err
	}
	rows, err := r.q.ListOrdersByUserFiltered(ctx, dbsqlc.ListOrdersByUserFilteredParams{
		UserID: userID, Column2: statuses, Limit: limit, Offset: offset,
	})
	return rows, total, err
}

func (r *pgOrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*dbsqlc.Order, error) {
	o, err := r.q.GetOrderByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *pgOrderRepository) Items(ctx context.Context, orderID uuid.UUID) ([]dbsqlc.ListOrderItemsWithProductRow, error) {
	return r.q.ListOrderItemsWithProduct(ctx, orderID)
}

func (r *pgOrderRepository) Address(ctx context.Context, orderID uuid.UUID) (*dbsqlc.OrderAddress, error) {
	a, err := r.q.GetOrderAddress(ctx, orderID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *pgOrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status dbsqlc.OrderStatus) (*dbsqlc.Order, error) {
	o, err := r.q.UpdateOrderStatus(ctx, dbsqlc.UpdateOrderStatusParams{ID: id, Status: status})
	if err != nil {
		return nil, err
	}
	return &o, nil
}
