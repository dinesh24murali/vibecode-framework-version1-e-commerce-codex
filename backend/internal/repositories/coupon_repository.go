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

// CouponRepository reads coupons and per-user redemptions.
type CouponRepository interface {
	// GetByCode returns an enabled coupon, or sql.ErrNoRows.
	GetByCode(ctx context.Context, code string) (*dbsqlc.Coupon, error)
	HasUsed(ctx context.Context, userID, couponID uuid.UUID) (bool, error)
}

type pgCouponRepository struct {
	q *dbsqlc.Queries
}

// NewCouponRepository returns a Postgres-backed CouponRepository.
func NewCouponRepository(pool *pgxpool.Pool) CouponRepository {
	return &pgCouponRepository{q: dbsqlc.New(stdlib.OpenDBFromPool(pool))}
}

func (r *pgCouponRepository) GetByCode(ctx context.Context, code string) (*dbsqlc.Coupon, error) {
	c, err := r.q.GetCouponByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *pgCouponRepository) HasUsed(ctx context.Context, userID, couponID uuid.UUID) (bool, error) {
	_, err := r.q.GetCouponUserByUserAndCoupon(ctx, dbsqlc.GetCouponUserByUserAndCouponParams{UserID: userID, CouponID: couponID})
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}
