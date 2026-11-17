package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
)

// AdminOrderRepository extends OrderRepository with the cross-customer listing used by admins.
type AdminOrderRepository interface {
	OrderRepository
	ListAll(ctx context.Context, statuses []dbsqlc.OrderStatus, limit, offset int32) ([]dbsqlc.Order, int64, error)
}

// NewAdminOrderRepository returns a Postgres-backed AdminOrderRepository.
func NewAdminOrderRepository(pool *pgxpool.Pool) AdminOrderRepository {
	return &pgOrderRepository{q: dbsqlc.New(stdlib.OpenDBFromPool(pool))}
}

func (r *pgOrderRepository) ListAll(ctx context.Context, statuses []dbsqlc.OrderStatus, limit, offset int32) ([]dbsqlc.Order, int64, error) {
	if statuses == nil {
		statuses = []dbsqlc.OrderStatus{}
	}
	total, err := r.q.AdminCountOrders(ctx, statuses)
	if err != nil {
		return nil, 0, err
	}
	rows, err := r.q.AdminListOrders(ctx, dbsqlc.AdminListOrdersParams{Statuses: statuses, PageLimit: limit, PageOffset: offset})
	return rows, total, err
}
