package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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

// AttributeAdminRepository is the write-side persistence interface for attribute definitions and values.
type AttributeAdminRepository interface {
	AttributeRepository
	ListPaged(ctx context.Context, limit, offset int32) ([]dbsqlc.ProductAttribute, int64, error)
	GetByCode(ctx context.Context, code string) (*dbsqlc.ProductAttribute, error)
	Create(ctx context.Context, arg dbsqlc.CreateAttributeParams) (*dbsqlc.ProductAttribute, error)
	Update(ctx context.Context, arg dbsqlc.UpdateAttributeParams) (*dbsqlc.ProductAttribute, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
	CreateValue(ctx context.Context, arg dbsqlc.CreateAttributeValueParams) (*dbsqlc.ProductAttributeValue, error)
	UpdateValue(ctx context.Context, arg dbsqlc.UpdateAttributeValueParams) (*dbsqlc.ProductAttributeValue, error)
	DeleteValue(ctx context.Context, id, attributeID uuid.UUID) (bool, error)
	CountAssignments(ctx context.Context, valueID uuid.UUID) (int64, error)
}

// IsUniqueViolation reports whether err is a Postgres unique-constraint violation.
func IsUniqueViolation(err error) bool {
	var pg *pgconn.PgError
	return errors.As(err, &pg) && pg.Code == "23505"
}

// NewAttributeAdminRepository returns a Postgres-backed AttributeAdminRepository.
func NewAttributeAdminRepository(pool *pgxpool.Pool) AttributeAdminRepository {
	return &pgAttributeRepository{q: dbsqlc.New(stdlib.OpenDBFromPool(pool))}
}

func (r *pgAttributeRepository) ListPaged(ctx context.Context, limit, offset int32) ([]dbsqlc.ProductAttribute, int64, error) {
	total, err := r.q.CountAttributes(ctx)
	if err != nil {
		return nil, 0, err
	}
	rows, err := r.q.ListAttributesPaged(ctx, dbsqlc.ListAttributesPagedParams{Limit: limit, Offset: offset})
	return rows, total, err
}

func (r *pgAttributeRepository) GetByCode(ctx context.Context, code string) (*dbsqlc.ProductAttribute, error) {
	a, err := r.q.GetAttributeByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *pgAttributeRepository) Create(ctx context.Context, arg dbsqlc.CreateAttributeParams) (*dbsqlc.ProductAttribute, error) {
	a, err := r.q.CreateAttribute(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *pgAttributeRepository) Update(ctx context.Context, arg dbsqlc.UpdateAttributeParams) (*dbsqlc.ProductAttribute, error) {
	a, err := r.q.UpdateAttribute(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *pgAttributeRepository) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	n, err := r.q.DeleteAttribute(ctx, id)
	return n > 0, err
}

func (r *pgAttributeRepository) CreateValue(ctx context.Context, arg dbsqlc.CreateAttributeValueParams) (*dbsqlc.ProductAttributeValue, error) {
	v, err := r.q.CreateAttributeValue(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *pgAttributeRepository) UpdateValue(ctx context.Context, arg dbsqlc.UpdateAttributeValueParams) (*dbsqlc.ProductAttributeValue, error) {
	v, err := r.q.UpdateAttributeValue(ctx, arg)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *pgAttributeRepository) DeleteValue(ctx context.Context, id, attributeID uuid.UUID) (bool, error) {
	n, err := r.q.DeleteAttributeValue(ctx, dbsqlc.DeleteAttributeValueParams{ID: id, AttributeID: attributeID})
	return n > 0, err
}

func (r *pgAttributeRepository) CountAssignments(ctx context.Context, valueID uuid.UUID) (int64, error) {
	return r.q.CountAssignmentsForValue(ctx, valueID)
}
