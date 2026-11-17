package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/models"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
)

// Errors surfaced by checkout.
var (
	ErrAddressNotFound   = errors.New("address not found")
	ErrInsufficientStock = repositories.ErrInsufficientStock
)

// Quote is the priced checkout summary.
type Quote struct { //nolint:govet // field order is for readability
	Items      []models.CartItem
	Address    *dbsqlc.Address
	CouponCode *string
	Subtotal   float64
	Discount   float64
	Taxable    float64
	GST        float64
	Total      float64
	draft      *repositories.OrderDraft
}

// CouponValidation is the result of validating a coupon on its own.
type CouponValidation struct { //nolint:govet // field order is for readability
	Valid   bool
	Code    string
	Type    *string
	Value   *float64
	Message string
}

// CheckoutService prices carts and places orders.
type CheckoutService struct {
	cart      *CartService
	coupons   *CouponService
	addresses repositories.AddressRepository
	orders    repositories.CheckoutRepository
}

// NewCheckoutService returns a CheckoutService.
func NewCheckoutService(cart *CartService, coupons *CouponService, addresses repositories.AddressRepository, orders repositories.CheckoutRepository) *CheckoutService {
	return &CheckoutService{cart: cart, coupons: coupons, addresses: addresses, orders: orders}
}

// ValidateCoupon reports whether the code can be used by the user. An unusable
// code is a normal outcome (is_valid=false), not an error.
func (s *CheckoutService) ValidateCoupon(ctx context.Context, userID uuid.UUID, code string) (*CouponValidation, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	app, err := s.coupons.Resolve(ctx, userID, code)
	if errors.Is(err, ErrCouponInvalid) {
		return &CouponValidation{Code: code, Message: CouponReason(err)}, nil
	}
	if err != nil {
		return nil, err
	}
	return &CouponValidation{Valid: true, Code: code, Type: &app.Type, Value: &app.Value, Message: app.Message}, nil
}

func (s *CheckoutService) quote(ctx context.Context, userID, addressID uuid.UUID, couponCode *string) (*Quote, error) {
	addr, err := s.addresses.GetForUser(ctx, addressID, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAddressNotFound
	}
	if err != nil {
		return nil, err
	}
	cart, err := s.cart.Cart(ctx, userID)
	if err != nil {
		return nil, err
	}
	var app *CouponApplication
	if couponCode != nil && strings.TrimSpace(*couponCode) != "" {
		if app, err = s.coupons.Resolve(ctx, userID, *couponCode); err != nil {
			return nil, err
		}
	}
	draft, err := NewOrderBuilder(userID).WithItems(cart.Items).WithAddress(addr).WithCoupon(app).Build()
	if err != nil {
		return nil, err
	}
	q := &Quote{
		Items: cart.Items, Address: addr, Subtotal: draft.Subtotal, Discount: draft.Discount,
		Taxable: draft.Taxable, GST: draft.GST, Total: draft.Total, draft: draft,
	}
	if draft.CouponCode != "" {
		q.CouponCode = &draft.CouponCode
	}
	return q, nil
}

// Calculate prices the user's cart for the address and optional coupon without creating an order.
func (s *CheckoutService) Calculate(ctx context.Context, userID, addressID uuid.UUID, couponCode *string) (*Quote, error) {
	return s.quote(ctx, userID, addressID, couponCode)
}

// PlaceOrder creates the order from the cart, reserving stock atomically.
func (s *CheckoutService) PlaceOrder(ctx context.Context, userID, addressID uuid.UUID, couponCode *string) (*OrderView, error) {
	q, err := s.quote(ctx, userID, addressID, couponCode)
	if err != nil {
		return nil, err
	}
	placed, err := s.orders.Place(ctx, q.draft)
	if errors.Is(err, repositories.ErrCouponAlreadyUsed) {
		return nil, couponErr("You have already used this coupon")
	}
	if err != nil {
		return nil, err
	}
	return &OrderView{Order: placed.Order, Address: placed.Address, Status: SpecStatus(placed.Order.Status), Items: placed.Items}, nil
}
