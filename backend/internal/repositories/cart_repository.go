package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
)

// CartRepository reads and writes a user's cart rows (one row per product).
type CartRepository interface {
	ListByUser(ctx context.Context, userID uuid.UUID) ([]dbsqlc.Cart, error)
	GetByProduct(ctx context.Context, userID, productID uuid.UUID) (*dbsqlc.Cart, error)
	Upsert(ctx context.Context, userID, productID uuid.UUID, quantity int32) error
	SetQuantity(ctx context.Context, id, userID uuid.UUID, quantity int32) error
	// Delete removes one row and reports how many rows were affected.
	Delete(ctx context.Context, id, userID uuid.UUID) (int64, error)
	Clear(ctx context.Context, userID uuid.UUID) error
}

type pgCartRepository struct {
	q *dbsqlc.Queries
}

// NewCartRepository returns a Postgres-backed CartRepository.
func NewCartRepository(pool *pgxpool.Pool) CartRepository {
	return &pgCartRepository{q: dbsqlc.New(stdlib.OpenDBFromPool(pool))}
}

func (r *pgCartRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]dbsqlc.Cart, error) {
	return r.q.GetCartByUser(ctx, userID)
}

func (r *pgCartRepository) GetByProduct(ctx context.Context, userID, productID uuid.UUID) (*dbsqlc.Cart, error) {
	row, err := r.q.GetCartItem(ctx, dbsqlc.GetCartItemParams{UserID: userID, ProductID: productID})
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *pgCartRepository) Upsert(ctx context.Context, userID, productID uuid.UUID, quantity int32) error {
	_, err := r.q.UpsertCartItem(ctx, dbsqlc.UpsertCartItemParams{UserID: userID, ProductID: productID, Quantity: quantity})
	return err
}

func (r *pgCartRepository) SetQuantity(ctx context.Context, id, userID uuid.UUID, quantity int32) error {
	_, err := r.q.UpdateCartItemQuantityByID(ctx, dbsqlc.UpdateCartItemQuantityByIDParams{ID: id, UserID: userID, Quantity: quantity})
	return err
}

func (r *pgCartRepository) Delete(ctx context.Context, id, userID uuid.UUID) (int64, error) {
	return r.q.DeleteCartItemByID(ctx, dbsqlc.DeleteCartItemByIDParams{ID: id, UserID: userID})
}

func (r *pgCartRepository) Clear(ctx context.Context, userID uuid.UUID) error {
	return r.q.ClearCart(ctx, userID)
}
