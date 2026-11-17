package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
)

// ContactRepository persists contact form enquiries.
type ContactRepository interface {
	Create(ctx context.Context, arg dbsqlc.CreateContactQueryParams) (dbsqlc.ContactQuery, error)
}

type pgContactRepository struct {
	q *dbsqlc.Queries
}

// NewContactRepository returns a Postgres-backed ContactRepository.
func NewContactRepository(pool *pgxpool.Pool) ContactRepository {
	return &pgContactRepository{q: dbsqlc.New(stdlib.OpenDBFromPool(pool))}
}

func (r *pgContactRepository) Create(ctx context.Context, arg dbsqlc.CreateContactQueryParams) (dbsqlc.ContactQuery, error) {
	return r.q.CreateContactQuery(ctx, arg)
}
