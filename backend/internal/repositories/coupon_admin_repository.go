package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
)

// CouponAdminRepository is the management persistence interface for coupons.
type CouponAdminRepository interface {
	ListPaged(ctx context.Context, limit, offset int32) ([]dbsqlc.Coupon, int64, error)
	Get(ctx context.Context, id uuid.UUID) (*dbsqlc.Coupon, error)
	Create(ctx context.Context, arg dbsqlc.CreateCouponParams) (*dbsqlc.Coupon, error)
	Update(ctx context.Context, arg dbsqlc.AdminUpdateCouponParams) (*dbsqlc.Coupon, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
	CountRedemptions(ctx context.Context, id uuid.UUID) (int64, error)
}

type pgCouponAdminRepository struct {
	q *dbsqlc.Queries
}

// NewCouponAdminRepository returns a Postgres-backed CouponAdminRepository.
func NewCouponAdminRepository(pool *pgxpool.Pool) CouponAdminRepository {
	return &pgCouponAdminRepository{q: dbsqlc.New(stdlib.OpenDBFromPool(pool))}
}

func (r *pgCouponAdminRepository) ListPaged(ctx context.Context, limit, offset int32) ([]dbsqlc.Coupon, int64, error) {
	total, err := r.q.CountCoupons(ctx)
	if err != nil {
		return nil, 0, err
	}
	rows, err := r.q.ListCoupons(ctx, dbsqlc.ListCouponsParams{Limit: limit, Offset: offset})
	return rows, total, err
}

func (r *pgCouponAdminRepository) Get(ctx context.Context, id uuid.UUID) (*dbsqlc.Coupon, error) {
	c, err := r.q.GetCouponByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *pgCouponAdminRepository) Create(ctx context.Context, arg dbsqlc.CreateCouponParams) (*dbsqlc.Coupon, error) {
	c, err := r.q.CreateCoupon(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *pgCouponAdminRepository) Update(ctx context.Context, arg dbsqlc.AdminUpdateCouponParams) (*dbsqlc.Coupon, error) {
	c, err := r.q.AdminUpdateCoupon(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *pgCouponAdminRepository) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	n, err := r.q.DeleteCoupon(ctx, id)
	return n > 0, err
}

func (r *pgCouponAdminRepository) CountRedemptions(ctx context.Context, id uuid.UUID) (int64, error) {
	return r.q.CountCouponRedemptions(ctx, id)
}
