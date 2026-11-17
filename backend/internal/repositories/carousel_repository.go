package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
)

// CarouselRepository reads the hero carousel slides.
type CarouselRepository interface {
	ListActive(ctx context.Context) ([]dbsqlc.CarouselSlide, error)
}

type pgCarouselRepository struct {
	q *dbsqlc.Queries
}

// NewCarouselRepository returns a Postgres-backed CarouselRepository.
func NewCarouselRepository(pool *pgxpool.Pool) CarouselRepository {
	return &pgCarouselRepository{q: dbsqlc.New(stdlib.OpenDBFromPool(pool))}
}

func (r *pgCarouselRepository) ListActive(ctx context.Context) ([]dbsqlc.CarouselSlide, error) {
	return r.q.ListActiveCarouselSlides(ctx)
}
