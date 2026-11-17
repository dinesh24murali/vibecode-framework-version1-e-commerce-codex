package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
)

// CategoryRepository is the read-only persistence interface for categories.
type CategoryRepository interface {
	List(ctx context.Context) ([]dbsqlc.Category, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dbsqlc.Category, error)
}

type pgCategoryRepository struct {
	q *dbsqlc.Queries
}

// NewCategoryRepository returns a Postgres-backed CategoryRepository.
func NewCategoryRepository(pool *pgxpool.Pool) CategoryRepository {
	return &pgCategoryRepository{q: dbsqlc.New(stdlib.OpenDBFromPool(pool))}
}

func (r *pgCategoryRepository) List(ctx context.Context) ([]dbsqlc.Category, error) {
	return r.q.ListCategories(ctx)
}

func (r *pgCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*dbsqlc.Category, error) {
	c, err := r.q.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
