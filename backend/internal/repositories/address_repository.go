package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
)

// AddressRepository is the persistence interface for saved addresses.
type AddressRepository interface {
	ListByUser(ctx context.Context, userID uuid.UUID) ([]dbsqlc.Address, error)
	GetForUser(ctx context.Context, id, userID uuid.UUID) (*dbsqlc.Address, error)
	// Create inserts an address; when IsDefault is true all other defaults are cleared atomically.
	Create(ctx context.Context, p dbsqlc.CreateAddressParams) (*dbsqlc.Address, error)
	Update(ctx context.Context, p dbsqlc.UpdateAddressParams) (*dbsqlc.Address, error)
	Delete(ctx context.Context, id, userID uuid.UUID) error
}

type pgAddressRepository struct {
	db *sql.DB
	q  *dbsqlc.Queries
}

// NewAddressRepository returns a Postgres-backed AddressRepository.
func NewAddressRepository(pool *pgxpool.Pool) AddressRepository {
	db := stdlib.OpenDBFromPool(pool)
	return &pgAddressRepository{db: db, q: dbsqlc.New(db)}
}

func (r *pgAddressRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]dbsqlc.Address, error) {
	return r.q.ListAddressesByUser(ctx, userID)
}

func (r *pgAddressRepository) GetForUser(ctx context.Context, id, userID uuid.UUID) (*dbsqlc.Address, error) {
	a, err := r.q.GetAddressForUser(ctx, dbsqlc.GetAddressForUserParams{ID: id, UserID: userID})
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *pgAddressRepository) Create(ctx context.Context, p dbsqlc.CreateAddressParams) (*dbsqlc.Address, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	q := r.q.WithTx(tx)
	if p.IsDefault {
		if cerr := q.ClearDefaultAddresses(ctx, p.UserID); cerr != nil {
			return nil, cerr
		}
	}
	a, err := q.CreateAddress(ctx, p)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *pgAddressRepository) Update(ctx context.Context, p dbsqlc.UpdateAddressParams) (*dbsqlc.Address, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	q := r.q.WithTx(tx)
	if p.IsDefault {
		if cerr := q.ClearDefaultAddresses(ctx, p.UserID); cerr != nil {
			return nil, cerr
		}
	}
	a, err := q.UpdateAddress(ctx, p)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *pgAddressRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return r.q.DeleteAddress(ctx, dbsqlc.DeleteAddressParams{ID: id, UserID: userID})
}
