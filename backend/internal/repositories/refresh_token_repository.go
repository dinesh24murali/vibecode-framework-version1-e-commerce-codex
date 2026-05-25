package repositories

import (
	"context"
	"net"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/sqlc-dev/pqtype"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
)

// RefreshTokenRepository is the persistence interface for refresh tokens.
type RefreshTokenRepository interface {
	Create(ctx context.Context, params dbsqlc.CreateRefreshTokenParams) (*dbsqlc.RefreshToken, error)
	GetByHash(ctx context.Context, hash string) (*dbsqlc.RefreshToken, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeFamily(ctx context.Context, familyID string) error
}

type pgRefreshTokenRepository struct {
	q *dbsqlc.Queries
}

// NewRefreshTokenRepository returns a Postgres-backed RefreshTokenRepository.
func NewRefreshTokenRepository(pool *pgxpool.Pool) RefreshTokenRepository {
	return &pgRefreshTokenRepository{q: dbsqlc.New(stdlib.OpenDBFromPool(pool))}
}

func (r *pgRefreshTokenRepository) Create(ctx context.Context, params dbsqlc.CreateRefreshTokenParams) (*dbsqlc.RefreshToken, error) {
	t, err := r.q.CreateRefreshToken(ctx, params)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *pgRefreshTokenRepository) GetByHash(ctx context.Context, hash string) (*dbsqlc.RefreshToken, error) {
	t, err := r.q.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *pgRefreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	return r.q.RevokeRefreshToken(ctx, id)
}

func (r *pgRefreshTokenRepository) RevokeFamily(ctx context.Context, familyID string) error {
	return r.q.RevokeTokenFamily(ctx, familyID)
}

// ParseInet converts a string IP address to pqtype.Inet for Postgres insertion.
// Falls back to 0.0.0.0 if the string cannot be parsed.
func ParseInet(ip string) pqtype.Inet {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		parsed = net.ParseIP("0.0.0.0")
	}
	return pqtype.Inet{IPNet: net.IPNet{IP: parsed, Mask: net.CIDRMask(32, 32)}, Valid: true}
}
