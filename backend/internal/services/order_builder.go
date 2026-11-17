package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/models"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
)

// GSTRate is the GST applied to the taxable amount (after the coupon deduction).
const GSTRate = 0.18

// ErrEmptyCart is returned when there is nothing to check out.
var ErrEmptyCart = errors.New("cart is empty")

// OrderBuilder assembles a priced OrderDraft step by step (Builder pattern).
// Pricing order: subtotal, minus coupon, gives the taxable amount, plus GST gives the total.
type OrderBuilder struct { //nolint:govet // field order is for readability
	draft  repositories.OrderDraft
	coupon *CouponApplication
}

// NewOrderBuilder starts a draft for userID.
func NewOrderBuilder(userID uuid.UUID) *OrderBuilder {
	return &OrderBuilder{draft: repositories.OrderDraft{UserID: userID}}
}

// WithItems sets the order lines from the (server-priced) cart lines.
func (b *OrderBuilder) WithItems(items []models.CartItem) *OrderBuilder {
	b.draft.Items = make([]repositories.DraftItem, 0, len(items))
	for _, i := range items {
		b.draft.Items = append(b.draft.Items, repositories.DraftItem{
			ProductID: i.Product.ID, Title: i.Product.Name, SKU: i.Product.Slug,
			Quantity: i.Quantity, UnitPrice: i.UnitPrice, LineTotal: i.LineTotal,
		})
	}
	return b
}

// WithAddress snapshots the shipping address into the draft.
func (b *OrderBuilder) WithAddress(a *dbsqlc.Address) *OrderBuilder {
	b.draft.Address = dbsqlc.CreateOrderAddressParams{
		RecipientName: a.RecipientName, Line1: a.Line1, Line2: a.Line2, City: a.City, State: a.State,
		PostalCode: a.PostalCode, CountryCode: a.CountryCode, PhoneE164: a.PhoneE164,
	}
	return b
}

// WithCoupon attaches a resolved coupon; nil means no coupon.
func (b *OrderBuilder) WithCoupon(c *CouponApplication) *OrderBuilder {
	b.coupon = c
	return b
}

// Build prices the draft.
func (b *OrderBuilder) Build() (*repositories.OrderDraft, error) {
	d := b.draft
	if len(d.Items) == 0 {
		return nil, ErrEmptyCart
	}
	for _, it := range d.Items {
		d.Subtotal += it.LineTotal
	}
	d.Subtotal = round2(d.Subtotal)
	if b.coupon != nil {
		d.Discount = b.coupon.DiscountFor(d.Subtotal)
		d.CouponCode = b.coupon.Coupon.Code
		d.CouponID = uuid.NullUUID{UUID: b.coupon.Coupon.ID, Valid: true}
	}
	d.Taxable = round2(d.Subtotal - d.Discount)
	d.GST = round2(d.Taxable * GSTRate)
	d.Total = round2(d.Taxable + d.GST)
	d.OrderNumber = newOrderNumber(time.Now())
	return &d, nil
}

func newOrderNumber(t time.Time) string {
	return fmt.Sprintf("ORD-%s-%s", t.UTC().Format("20060102"), strings.ToUpper(uuid.NewString()[:6]))
}
