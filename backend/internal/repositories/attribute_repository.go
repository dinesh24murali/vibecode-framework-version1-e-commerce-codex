package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
)

// AttributeRepository is the read-only persistence interface for product attributes.
type AttributeRepository interface {
	ListFilterable(ctx context.Context) ([]dbsqlc.ProductAttribute, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dbsqlc.ProductAttribute, error)
	ListValues(ctx context.Context, attributeID uuid.UUID) ([]dbsqlc.ProductAttributeValue, error)
}

type pgAttributeRepository struct {
	q *dbsqlc.Queries
}

// NewAttributeRepository returns a Postgres-backed AttributeRepository.
func NewAttributeRepository(pool *pgxpool.Pool) AttributeRepository {
	return &pgAttributeRepository{q: dbsqlc.New(stdlib.OpenDBFromPool(pool))}
}

func (r *pgAttributeRepository) ListFilterable(ctx context.Context) ([]dbsqlc.ProductAttribute, error) {
	return r.q.ListFilterableAttributes(ctx)
}

func (r *pgAttributeRepository) GetByID(ctx context.Context, id uuid.UUID) (*dbsqlc.ProductAttribute, error) {
	a, err := r.q.GetAttributeByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *pgAttributeRepository) ListValues(ctx context.Context, attributeID uuid.UUID) ([]dbsqlc.ProductAttributeValue, error) {
	return r.q.ListValuesByAttribute(ctx, attributeID)
}
