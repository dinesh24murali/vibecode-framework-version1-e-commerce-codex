package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
)

// CategoryAdminRepository is the write-side persistence interface for categories.
type CategoryAdminRepository interface {
	CategoryRepository
	ListPaged(ctx context.Context, limit, offset int32) ([]dbsqlc.Category, int64, error)
	Create(ctx context.Context, slug, name string, imageURL *string) (*dbsqlc.Category, error)
	Update(ctx context.Context, id uuid.UUID, slug, name string, imageURL *string) (*dbsqlc.Category, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
	CountProducts(ctx context.Context, id uuid.UUID) (int64, error)
}

// NewCategoryAdminRepository returns a Postgres-backed CategoryAdminRepository.
func NewCategoryAdminRepository(pool *pgxpool.Pool) CategoryAdminRepository {
	return &pgCategoryRepository{q: dbsqlc.New(stdlib.OpenDBFromPool(pool))}
}

// IsForeignKeyViolation reports whether err is a Postgres foreign-key violation.
func IsForeignKeyViolation(err error) bool {
	var pg *pgconn.PgError
	return errors.As(err, &pg) && pg.Code == "23503"
}

func nullStr(s *string) sql.NullString {
	if s == nil || *s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

func (r *pgCategoryRepository) ListPaged(ctx context.Context, limit, offset int32) ([]dbsqlc.Category, int64, error) {
	total, err := r.q.CountCategories(ctx)
	if err != nil {
		return nil, 0, err
	}
	rows, err := r.q.ListCategoriesPaged(ctx, dbsqlc.ListCategoriesPagedParams{Limit: limit, Offset: offset})
	return rows, total, err
}

func (r *pgCategoryRepository) Create(ctx context.Context, slug, name string, imageURL *string) (*dbsqlc.Category, error) {
	c, err := r.q.AdminCreateCategory(ctx, dbsqlc.AdminCreateCategoryParams{Slug: slug, Name: name, ImageUrl: nullStr(imageURL)})
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *pgCategoryRepository) Update(ctx context.Context, id uuid.UUID, slug, name string, imageURL *string) (*dbsqlc.Category, error) {
	c, err := r.q.AdminUpdateCategory(ctx, dbsqlc.AdminUpdateCategoryParams{ID: id, Slug: slug, Name: name, ImageUrl: nullStr(imageURL)})
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *pgCategoryRepository) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	n, err := r.q.AdminDeleteCategory(ctx, id)
	return n > 0, err
}

func (r *pgCategoryRepository) CountProducts(ctx context.Context, id uuid.UUID) (int64, error) {
	return r.q.CountProductsInCategory(ctx, id)
}
