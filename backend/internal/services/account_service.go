package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
)

// Sentinel errors returned by AccountService.
var (
	ErrNotFound         = errors.New("not found")
	ErrWrongPassword    = errors.New("current password is incorrect")
	ErrPasswordMismatch = errors.New("new password and confirmation do not match")
	ErrOrderNotCancel   = errors.New("order cannot be cancelled in its current status")
	ErrInvalidStatus    = errors.New("invalid order status filter")
)

// OTPPurposeEmailChange is the OTP purpose used when changing the account email.
const OTPPurposeEmailChange = "email_change_otp"

// OTPStore is the subset of the OTP service the account flows need.
type OTPStore interface {
	GenerateAndStore(ctx context.Context, email, purpose string) (string, error)
	Verify(ctx context.Context, email, purpose, otp string) (bool, error)
}

// AddressInput carries the writable Address fields from the API spec.
type AddressInput struct {
	AddressLine1 string
	AddressLine2 *string
	Landmark     *string
	City         string
	State        string
	Pincode      string
	CompanyName  *string
	GSTNumber    *string
	Phone        string
	IsDefault    bool
}

// OrderPage is one page of orders plus the total count.
type OrderPage struct {
	Orders []OrderView
	Total  int64
}

// OrderView is an order with its items and shipping address, in spec terms.
type OrderView struct {
	Order   *dbsqlc.Order
	Address *dbsqlc.OrderAddress
	// Status is the spec status: open, fulfilled or cancelled.
	Status string
	Items  []dbsqlc.ListOrderItemsWithProductRow
}

// AccountService implements profile, address and order-history flows.
type AccountService struct {
	users     repositories.UserRepository
	addresses repositories.AddressRepository
	orders    repositories.OrderRepository
	otp       OTPStore
}

// NewAccountService constructs an AccountService.
func NewAccountService(
	users repositories.UserRepository,
	addresses repositories.AddressRepository,
	orders repositories.OrderRepository,
	otp OTPStore,
) *AccountService {
	return &AccountService{users: users, addresses: addresses, orders: orders, otp: otp}
}

// ── Profile ──────────────────────────────────────────────────────────────────

// GetProfile returns the user.
func (s *AccountService) GetProfile(ctx context.Context, userID uuid.UUID) (*dbsqlc.User, error) {
	u, err := s.users.GetByID(ctx, userID)
	return u, mapNotFound(err)
}

// UpdateProfile stores first and last name (users.full_name is a single column).
func (s *AccountService) UpdateProfile(ctx context.Context, userID uuid.UUID, first, last string) (*dbsqlc.User, error) {
	full := strings.TrimSpace(strings.TrimSpace(first) + " " + strings.TrimSpace(last))
	u, err := s.users.UpdateName(ctx, userID, full)
	return u, mapNotFound(err)
}

// ChangePassword verifies the current password and sets a new one.
func (s *AccountService) ChangePassword(ctx context.Context, userID uuid.UUID, current, newPw, confirm string) error {
	if newPw != confirm {
		return ErrPasswordMismatch
	}
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return mapNotFound(err)
	}
	if !u.PasswordHash.Valid || !verifyPassword(current, u.PasswordHash.String) {
		return ErrWrongPassword
	}
	hash, err := hashPassword(newPw)
	if err != nil {
		return err
	}
	return s.users.UpdatePasswordHash(ctx, userID, hash)
}

// RequestEmailChange sends an OTP to the new address. It always succeeds from the
// caller's point of view to avoid revealing which emails are registered.
func (s *AccountService) RequestEmailChange(ctx context.Context, newEmail string) {
	newEmail = strings.ToLower(strings.TrimSpace(newEmail))
	if existing, err := s.users.GetByEmail(ctx, newEmail); err == nil && existing != nil {
		return
	}
	code, err := s.otp.GenerateAndStore(ctx, newEmail, OTPPurposeEmailChange)
	if err != nil {
		slog.Error("email change: store OTP", "error", err)
		return
	}
	// DEV: log OTP to stdout. Replace with SMTP/SES in production.
	slog.Info("email change OTP", "email", newEmail, "otp", code)
}

// ConfirmEmailChange verifies the OTP and commits the new email.
func (s *AccountService) ConfirmEmailChange(ctx context.Context, userID uuid.UUID, newEmail, code string) (*dbsqlc.User, error) {
	newEmail = strings.ToLower(strings.TrimSpace(newEmail))
	ok, err := s.otp.Verify(ctx, newEmail, OTPPurposeEmailChange, code)
	if err != nil {
		return nil, fmt.Errorf("confirm email change: %w", err)
	}
	if !ok {
		return nil, ErrInvalidOTP
	}
	if existing, gerr := s.users.GetByEmail(ctx, newEmail); gerr == nil && existing != nil && existing.ID != userID {
		return nil, ErrEmailAlreadyExists
	}
	u, err := s.users.UpdateEmail(ctx, userID, newEmail)
	return u, mapNotFound(err)
}

// ── Addresses ────────────────────────────────────────────────────────────────

// ListAddresses returns the user's addresses, default first.
func (s *AccountService) ListAddresses(ctx context.Context, userID uuid.UUID) ([]dbsqlc.Address, error) {
	return s.addresses.ListByUser(ctx, userID)
}

// CreateAddress saves a new address.
func (s *AccountService) CreateAddress(ctx context.Context, userID uuid.UUID, in AddressInput) (*dbsqlc.Address, error) {
	return s.addresses.Create(ctx, dbsqlc.CreateAddressParams{
		UserID: userID, Line1: in.AddressLine1, Line2: nullStr(in.AddressLine2),
		Landmark: nullStr(in.Landmark), City: in.City, State: in.State, PostalCode: in.Pincode,
		CompanyName: nullStr(in.CompanyName), GstNumber: nullStr(in.GSTNumber),
		PhoneE164: sql.NullString{String: in.Phone, Valid: in.Phone != ""}, IsDefault: in.IsDefault,
	})
}

// UpdateAddress replaces an existing address owned by the user.
func (s *AccountService) UpdateAddress(ctx context.Context, userID, id uuid.UUID, in AddressInput) (*dbsqlc.Address, error) {
	a, err := s.addresses.Update(ctx, dbsqlc.UpdateAddressParams{
		ID: id, UserID: userID, Line1: in.AddressLine1, Line2: nullStr(in.AddressLine2),
		Landmark: nullStr(in.Landmark), City: in.City, State: in.State, PostalCode: in.Pincode,
		CompanyName: nullStr(in.CompanyName), GstNumber: nullStr(in.GSTNumber),
		PhoneE164: sql.NullString{String: in.Phone, Valid: in.Phone != ""}, IsDefault: in.IsDefault,
	})
	return a, mapNotFound(err)
}

// DeleteAddress removes an address owned by the user.
func (s *AccountService) DeleteAddress(ctx context.Context, userID, id uuid.UUID) error {
	if _, err := s.addresses.GetForUser(ctx, id, userID); err != nil {
		return mapNotFound(err)
	}
	return s.addresses.Delete(ctx, id, userID)
}

// ── Orders ───────────────────────────────────────────────────────────────────

// SpecStatus maps a DB order status to the API statuses open/fulfilled/cancelled.
func SpecStatus(st dbsqlc.OrderStatus) string {
	switch st {
	case dbsqlc.OrderStatusDelivered:
		return "fulfilled"
	case dbsqlc.OrderStatusCancelled, dbsqlc.OrderStatusPaymentFailed, dbsqlc.OrderStatusRefunded:
		return "cancelled"
	default:
		return "open"
	}
}

func dbStatuses(spec string) ([]dbsqlc.OrderStatus, error) {
	switch spec {
	case "":
		return nil, nil
	case "open":
		return []dbsqlc.OrderStatus{dbsqlc.OrderStatusPendingPayment, dbsqlc.OrderStatusPaid, dbsqlc.OrderStatusPacked, dbsqlc.OrderStatusShipped}, nil
	case "fulfilled":
		return []dbsqlc.OrderStatus{dbsqlc.OrderStatusDelivered}, nil
	case "cancelled":
		return []dbsqlc.OrderStatus{dbsqlc.OrderStatusCancelled, dbsqlc.OrderStatusPaymentFailed, dbsqlc.OrderStatusRefunded}, nil
	}
	return nil, ErrInvalidStatus
}

// ListOrders returns one page of the user's orders, optionally filtered by spec status.
func (s *AccountService) ListOrders(ctx context.Context, userID uuid.UUID, status string, page, limit int) (*OrderPage, error) {
	statuses, err := dbStatuses(status)
	if err != nil {
		return nil, err
	}
	rows, total, err := s.orders.ListByUser(ctx, userID, statuses, int32(limit), int32((page-1)*limit)) //nolint:gosec // bounded by handler
	if err != nil {
		return nil, err
	}
	out := &OrderPage{Total: total, Orders: make([]OrderView, 0, len(rows))}
	for i := range rows {
		v, err := s.view(ctx, &rows[i])
		if err != nil {
			return nil, err
		}
		out.Orders = append(out.Orders, *v)
	}
	return out, nil
}

// GetOrder returns one of the user's orders. Other users' orders yield ErrForbidden.
func (s *AccountService) GetOrder(ctx context.Context, userID, id uuid.UUID) (*OrderView, error) {
	o, err := s.ownedOrder(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	return s.view(ctx, o)
}

// CancelOrder cancels an order that has not yet shipped.
func (s *AccountService) CancelOrder(ctx context.Context, userID, id uuid.UUID) (*OrderView, error) {
	o, err := s.ownedOrder(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	switch o.Status {
	case dbsqlc.OrderStatusPendingPayment, dbsqlc.OrderStatusPaid, dbsqlc.OrderStatusPacked:
	default:
		return nil, ErrOrderNotCancel
	}
	updated, err := s.orders.UpdateStatus(ctx, id, dbsqlc.OrderStatusCancelled)
	if err != nil {
		return nil, err
	}
	return s.view(ctx, updated)
}

// ErrForbidden is returned when an order belongs to another user.
var ErrForbidden = errors.New("forbidden")

func (s *AccountService) ownedOrder(ctx context.Context, userID, id uuid.UUID) (*dbsqlc.Order, error) {
	o, err := s.orders.GetByID(ctx, id)
	if err != nil {
		return nil, mapNotFound(err)
	}
	if o.UserID != userID {
		return nil, ErrForbidden
	}
	return o, nil
}

func (s *AccountService) view(ctx context.Context, o *dbsqlc.Order) (*OrderView, error) {
	items, err := s.orders.Items(ctx, o.ID)
	if err != nil {
		return nil, err
	}
	addr, err := s.orders.Address(ctx, o.ID)
	if err != nil {
		return nil, err
	}
	return &OrderView{Order: o, Items: items, Address: addr, Status: SpecStatus(o.Status)}, nil
}

// ── helpers ──────────────────────────────────────────────────────────────────

func nullStr(p *string) sql.NullString {
	if p == nil || *p == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: *p, Valid: true}
}

func mapNotFound(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}
