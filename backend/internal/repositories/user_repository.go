// Package repositories provides interfaces and Postgres implementations for
// all data access. Handlers and services depend on the interfaces, not the
// concrete types, enabling clean testing boundaries.
package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
)

// UserRepository is the persistence interface for the users domain.
type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*dbsqlc.User, error)
	GetByEmail(ctx context.Context, email string) (*dbsqlc.User, error)
	Create(ctx context.Context, params dbsqlc.CreateUserParams) (*dbsqlc.User, error)
	UpdatePasswordHash(ctx context.Context, id uuid.UUID, hash string) error
}

type pgUserRepository struct {
	q *dbsqlc.Queries
}

// NewUserRepository returns a Postgres-backed UserRepository.
func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &pgUserRepository{q: dbsqlc.New(stdlib.OpenDBFromPool(pool))}
}

func (r *pgUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*dbsqlc.User, error) {
	u, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *pgUserRepository) GetByEmail(ctx context.Context, email string) (*dbsqlc.User, error) {
	u, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *pgUserRepository) Create(ctx context.Context, params dbsqlc.CreateUserParams) (*dbsqlc.User, error) {
	u, err := r.q.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *pgUserRepository) UpdatePasswordHash(ctx context.Context, id uuid.UUID, hash string) error {
	return r.q.UpdatePasswordHash(ctx, dbsqlc.UpdatePasswordHashParams{
		ID:           id,
		PasswordHash: sql.NullString{String: hash, Valid: true},
	})
}
