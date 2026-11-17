package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
)

// ErrCouponInvalid wraps the reason a coupon cannot be applied.
var ErrCouponInvalid = errors.New("coupon invalid")

// CouponApplication is a coupon resolved against a subtotal.
type CouponApplication struct { //nolint:govet // field order is for readability
	Coupon *dbsqlc.Coupon
	// Type is "flat" or "percentage"; Value is the INR amount or the percent.
	Type    string
	Value   float64
	Message string
}

// DiscountFor returns the discount for subtotal, never more than the subtotal.
func (a *CouponApplication) DiscountFor(subtotal float64) float64 {
	d := a.Value
	if a.Type == "percentage" {
		d = subtotal * a.Value / 100
	}
	if d > subtotal {
		d = subtotal
	}
	return round2(d)
}

// CouponService validates coupon codes.
type CouponService struct {
	repo repositories.CouponRepository
	now  func() time.Time
}

// NewCouponService returns a CouponService.
func NewCouponService(repo repositories.CouponRepository) *CouponService {
	return &CouponService{repo: repo, now: time.Now}
}

func couponErr(msg string) error { return fmt.Errorf("%w: %s", ErrCouponInvalid, msg) }

// CouponReason extracts the user-facing reason from an ErrCouponInvalid error.
func CouponReason(err error) string {
	return strings.TrimPrefix(err.Error(), ErrCouponInvalid.Error()+": ")
}

// Resolve checks that code is enabled, unexpired and not yet used by the user.
func (s *CouponService) Resolve(ctx context.Context, userID uuid.UUID, code string) (*CouponApplication, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" {
		return nil, couponErr("Enter a coupon code")
	}
	c, err := s.repo.GetByCode(ctx, code)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, couponErr("Invalid coupon code")
	}
	if err != nil {
		return nil, err
	}
	if c.CouponExpiry.Valid && c.CouponExpiry.Time.Before(s.now()) {
		return nil, couponErr("This coupon has expired")
	}
	used, err := s.repo.HasUsed(ctx, userID, c.ID)
	if err != nil {
		return nil, err
	}
	if used {
		return nil, couponErr("You have already used this coupon")
	}
	app := &CouponApplication{Coupon: c, Type: c.Type}
	if c.Type == "percentage" {
		app.Value = float64(c.Percent.Int32)
		app.Message = fmt.Sprintf("Coupon applied — %g%% off your order", app.Value)
	} else {
		app.Value = float64(c.Flat.Int32)
		app.Message = fmt.Sprintf("Coupon applied — INR %g off your order", app.Value)
	}
	return app, nil
}
