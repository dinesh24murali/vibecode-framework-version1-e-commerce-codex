package services

import (
	"context"
	"errors"

	"github.com/google/uuid"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
)

// ErrInvalidTransition is returned when an order cannot move to the requested status.
var ErrInvalidTransition = errors.New("invalid order status transition")

// AdminOrderView is an order plus the customer it belongs to.
type AdminOrderView struct { //nolint:govet // embedding order matters more than packing
	OrderView
	CustomerName  string
	CustomerEmail string
}

// AdminOrderPage is one page of admin order views plus the total.
type AdminOrderPage struct {
	Orders []AdminOrderView
	Total  int64
}

// AdminOrderService implements admin order management.
type AdminOrderService struct {
	orders repositories.AdminOrderRepository
	users  repositories.UserRepository
}

// NewAdminOrderService constructs an AdminOrderService.
func NewAdminOrderService(orders repositories.AdminOrderRepository, users repositories.UserRepository) *AdminOrderService {
	return &AdminOrderService{orders: orders, users: users}
}

func (s *AdminOrderService) view(ctx context.Context, o *dbsqlc.Order) (*AdminOrderView, error) {
	items, err := s.orders.Items(ctx, o.ID)
	if err != nil {
		return nil, err
	}
	addr, err := s.orders.Address(ctx, o.ID)
	if err != nil {
		return nil, err
	}
	v := &AdminOrderView{OrderView: OrderView{Order: o, Items: items, Address: addr, Status: SpecStatus(o.Status)}}
	if u, err := s.users.GetByID(ctx, o.UserID); err == nil {
		v.CustomerName, v.CustomerEmail = u.FullName, u.Email
	}
	if v.CustomerName == "" && addr != nil {
		v.CustomerName = addr.RecipientName
	}
	return v, nil
}

// List returns one page of all orders, optionally filtered by spec status (open/fulfilled/cancelled).
func (s *AdminOrderService) List(ctx context.Context, status string, page, limit int) (*AdminOrderPage, error) {
	statuses, err := dbStatuses(status)
	if err != nil {
		return nil, err
	}
	rows, total, err := s.orders.ListAll(ctx, statuses, int32(limit), int32((page-1)*limit)) //nolint:gosec // bounded by handler
	if err != nil {
		return nil, err
	}
	out := &AdminOrderPage{Total: total, Orders: make([]AdminOrderView, 0, len(rows))}
	for i := range rows {
		v, err := s.view(ctx, &rows[i])
		if err != nil {
			return nil, err
		}
		out.Orders = append(out.Orders, *v)
	}
	return out, nil
}

// Get returns one order or ErrNotFound.
func (s *AdminOrderService) Get(ctx context.Context, id uuid.UUID) (*AdminOrderView, error) {
	o, err := s.orders.GetByID(ctx, id)
	if err != nil {
		return nil, mapNotFound(err)
	}
	return s.view(ctx, o)
}

// UpdateStatus moves an order to "fulfilled" (delivered) or "cancelled".
// Only orders still open can change; unpaid orders cannot be fulfilled. Any other
// request (including "open") is ErrInvalidTransition.
func (s *AdminOrderService) UpdateStatus(ctx context.Context, id uuid.UUID, target string) (*AdminOrderView, error) {
	o, err := s.orders.GetByID(ctx, id)
	if err != nil {
		return nil, mapNotFound(err)
	}
	if SpecStatus(o.Status) != "open" {
		return nil, ErrInvalidTransition
	}
	var next dbsqlc.OrderStatus
	switch target {
	case "fulfilled":
		if o.Status == dbsqlc.OrderStatusPendingPayment {
			return nil, ErrInvalidTransition
		}
		next = dbsqlc.OrderStatusDelivered
	case "cancelled":
		next = dbsqlc.OrderStatusCancelled
	default:
		return nil, ErrInvalidTransition
	}
	updated, err := s.orders.UpdateStatus(ctx, id, next)
	if err != nil {
		return nil, mapNotFound(err)
	}
	return s.view(ctx, updated)
}
