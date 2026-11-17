package services

import (
	"context"
	"database/sql"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
)

var couponCodeRe = regexp.MustCompile(`^[A-Z0-9][A-Z0-9_-]{1,39}$`)

// NullableTime distinguishes "absent" from an explicit null/value in update requests.
type NullableTime struct {
	Time  time.Time
	Set   bool
	Valid bool
}

// AdminCouponInput carries create/update fields; nil/unset means "not provided".
type AdminCouponInput struct {
	Name        *string
	Description *string
	Code        *string
	IsEnable    *bool
	Type        *string
	Flat        *int32
	Percent     *int32
	Expiry      NullableTime
}

// AdminCouponService implements admin coupon management.
type AdminCouponService struct {
	repo repositories.CouponAdminRepository
}

// NewAdminCouponService constructs an AdminCouponService.
func NewAdminCouponService(repo repositories.CouponAdminRepository) *AdminCouponService {
	return &AdminCouponService{repo: repo}
}

// List returns one page of coupons and the total.
func (s *AdminCouponService) List(ctx context.Context, page, limit int) ([]dbsqlc.Coupon, int64, error) {
	return s.repo.ListPaged(ctx, int32(limit), int32((page-1)*limit)) //nolint:gosec // bounded by handler
}

// Get returns one coupon or ErrNotFound.
func (s *AdminCouponService) Get(ctx context.Context, id uuid.UUID) (*dbsqlc.Coupon, error) {
	c, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, mapNotFound(err)
	}
	return c, nil
}

// amounts validates type against flat/percent and returns the column values:
// only the value matching the type is stored, the other is NULL.
func amounts(typ string, flat, percent *int32) (sql.NullInt32, sql.NullInt32, error) {
	switch typ {
	case "flat":
		if flat == nil || *flat < 1 {
			return sql.NullInt32{}, sql.NullInt32{}, ErrValidation
		}
		return sql.NullInt32{}, sql.NullInt32{Int32: *flat, Valid: true}, nil
	case "percentage":
		if percent == nil || *percent < 1 || *percent > 100 {
			return sql.NullInt32{}, sql.NullInt32{}, ErrValidation
		}
		return sql.NullInt32{Int32: *percent, Valid: true}, sql.NullInt32{}, nil
	}
	return sql.NullInt32{}, sql.NullInt32{}, ErrValidation
}

func nullText(s *string) sql.NullString {
	if s == nil || strings.TrimSpace(*s) == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: strings.TrimSpace(*s), Valid: true}
}

// Create validates and creates a coupon; a duplicate code yields ErrConflict.
func (s *AdminCouponService) Create(ctx context.Context, in AdminCouponInput) (*dbsqlc.Coupon, error) {
	if in.Name == nil || in.Code == nil || in.Type == nil || strings.TrimSpace(*in.Name) == "" {
		return nil, ErrValidation
	}
	code := strings.ToUpper(strings.TrimSpace(*in.Code))
	if !couponCodeRe.MatchString(code) {
		return nil, ErrValidation
	}
	percent, flat, err := amounts(*in.Type, in.Flat, in.Percent)
	if err != nil {
		return nil, err
	}
	enabled := true
	if in.IsEnable != nil {
		enabled = *in.IsEnable
	}
	c, err := s.repo.Create(ctx, dbsqlc.CreateCouponParams{
		Name: strings.TrimSpace(*in.Name), Description: nullText(in.Description), Code: code, IsEnable: enabled,
		Type: *in.Type, Percent: percent, Flat: flat,
		CouponExpiry: sql.NullTime{Time: in.Expiry.Time, Valid: in.Expiry.Set && in.Expiry.Valid},
	})
	if err != nil {
		if repositories.IsUniqueViolation(err) {
			return nil, ErrConflict
		}
		return nil, err
	}
	return c, nil
}

// Update applies a partial update (the code is immutable).
func (s *AdminCouponService) Update(ctx context.Context, id uuid.UUID, in AdminCouponInput) (*dbsqlc.Coupon, error) {
	cur, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, mapNotFound(err)
	}
	arg := dbsqlc.AdminUpdateCouponParams{
		ID: id, Name: cur.Name, Description: cur.Description, IsEnable: cur.IsEnable,
		Type: cur.Type, Percent: cur.Percent, Flat: cur.Flat, CouponExpiry: cur.CouponExpiry,
	}
	if in.Name != nil {
		if arg.Name = strings.TrimSpace(*in.Name); arg.Name == "" {
			return nil, ErrValidation
		}
	}
	if in.Description != nil {
		arg.Description = nullText(in.Description)
	}
	if in.IsEnable != nil {
		arg.IsEnable = *in.IsEnable
	}
	if in.Type != nil || in.Flat != nil || in.Percent != nil {
		typ := cur.Type
		if in.Type != nil {
			typ = *in.Type
		}
		flat, percent := in.Flat, in.Percent
		if flat == nil && cur.Flat.Valid {
			flat = &cur.Flat.Int32
		}
		if percent == nil && cur.Percent.Valid {
			percent = &cur.Percent.Int32
		}
		arg.Type = typ
		if arg.Percent, arg.Flat, err = amounts(typ, flat, percent); err != nil {
			return nil, err
		}
	}
	if in.Expiry.Set {
		arg.CouponExpiry = sql.NullTime{Time: in.Expiry.Time, Valid: in.Expiry.Valid}
	}
	c, err := s.repo.Update(ctx, arg)
	if err != nil {
		return nil, mapNotFound(err)
	}
	return c, nil
}

// Delete removes a coupon; a redeemed coupon yields ErrConflict (disable it instead).
func (s *AdminCouponService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.repo.Get(ctx, id); err != nil {
		return mapNotFound(err)
	}
	n, err := s.repo.CountRedemptions(ctx, id)
	if err != nil {
		return err
	}
	if n > 0 {
		return ErrConflict
	}
	ok, err := s.repo.Delete(ctx, id)
	if err != nil {
		if repositories.IsForeignKeyViolation(err) {
			return ErrConflict
		}
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}
