package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
)

// CatalogFilter holds the optional public catalogue filters.
type CatalogFilter struct {
	Query         string
	Sort          string
	CategorySlugs []string
	ValueIDs      []uuid.UUID
}

// ProductRepository is the read-only persistence interface for the public catalogue.
// Only active products (is_active = true) are ever returned.
type ProductRepository interface {
	Count(ctx context.Context, f CatalogFilter) (int64, error)
	List(ctx context.Context, f CatalogFilter, limit, offset int32) ([]dbsqlc.ListCatalogProductsRow, error)
	GetActive(ctx context.Context, id uuid.UUID) (*dbsqlc.GetActiveProductByIDRow, error)
	Images(ctx context.Context, productIDs []uuid.UUID) ([]dbsqlc.ProductImage, error)
	Attributes(ctx context.Context, productIDs []uuid.UUID) ([]dbsqlc.ListAttributesByProductsRow, error)
	Inventory(ctx context.Context, productIDs []uuid.UUID) ([]dbsqlc.InventoryItem, error)
	ActiveCategorySlugs(ctx context.Context) ([]string, error)
}

type pgProductRepository struct {
	q *dbsqlc.Queries
}

// NewProductRepository returns a Postgres-backed ProductRepository.
func NewProductRepository(pool *pgxpool.Pool) ProductRepository {
	return &pgProductRepository{q: dbsqlc.New(stdlib.OpenDBFromPool(pool))}
}

func nonNil[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}

func nullQuery(q string) sql.NullString {
	return sql.NullString{String: q, Valid: q != ""}
}

func (r *pgProductRepository) Count(ctx context.Context, f CatalogFilter) (int64, error) {
	return r.q.CountCatalogProducts(ctx, dbsqlc.CountCatalogProductsParams{
		CategorySlugs: nonNil(f.CategorySlugs), Query: nullQuery(f.Query), ValueIds: nonNil(f.ValueIDs),
	})
}

func (r *pgProductRepository) List(ctx context.Context, f CatalogFilter, limit, offset int32) ([]dbsqlc.ListCatalogProductsRow, error) {
	return r.q.ListCatalogProducts(ctx, dbsqlc.ListCatalogProductsParams{
		CategorySlugs: nonNil(f.CategorySlugs), Query: nullQuery(f.Query), ValueIds: nonNil(f.ValueIDs),
		Sort: f.Sort, PageLimit: limit, PageOffset: offset,
	})
}

func (r *pgProductRepository) GetActive(ctx context.Context, id uuid.UUID) (*dbsqlc.GetActiveProductByIDRow, error) {
	p, err := r.q.GetActiveProductByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *pgProductRepository) Images(ctx context.Context, ids []uuid.UUID) ([]dbsqlc.ProductImage, error) {
	return r.q.ListImagesByProducts(ctx, nonNil(ids))
}

func (r *pgProductRepository) Attributes(ctx context.Context, ids []uuid.UUID) ([]dbsqlc.ListAttributesByProductsRow, error) {
	return r.q.ListAttributesByProducts(ctx, nonNil(ids))
}

func (r *pgProductRepository) Inventory(ctx context.Context, ids []uuid.UUID) ([]dbsqlc.InventoryItem, error) {
	return r.q.ListInventoryByProducts(ctx, nonNil(ids))
}

func (r *pgProductRepository) ActiveCategorySlugs(ctx context.Context) ([]string, error) {
	return r.q.ListActiveCategorySlugs(ctx)
}
