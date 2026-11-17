package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

// CheckoutHandler serves the authenticated checkout endpoints.
type CheckoutHandler struct {
	svc *services.CheckoutService
}

// NewCheckoutHandler returns a CheckoutHandler.
func NewCheckoutHandler(svc *services.CheckoutService) *CheckoutHandler {
	return &CheckoutHandler{svc: svc}
}

type checkoutRequest struct { //nolint:govet // field order is for readability
	AddressID  uuid.UUID `json:"address_id" binding:"required"`
	CouponCode *string   `json:"coupon_code"`
}

type couponRequest struct {
	CouponCode string `json:"coupon_code" binding:"required"`
}

func (h *CheckoutHandler) fail(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrEmptyCart):
		apiError(c, http.StatusBadRequest, "cart is empty")
	case errors.Is(err, services.ErrCouponInvalid):
		apiError(c, http.StatusUnprocessableEntity, services.CouponReason(err))
	case errors.Is(err, services.ErrAddressNotFound):
		apiError(c, http.StatusUnprocessableEntity, "address not found")
	case errors.Is(err, services.ErrInsufficientStock):
		apiError(c, http.StatusUnprocessableEntity, err.Error())
	default:
		apiError(c, http.StatusInternalServerError, "internal error")
	}
}

// Calculate godoc POST /api/v1/checkout/calculate
func (h *CheckoutHandler) Calculate(c *gin.Context) {
	var req checkoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, "address_id is required")
		return
	}
	q, err := h.svc.Calculate(c.Request.Context(), currentUser(c), req.AddressID, req.CouponCode)
	if err != nil {
		h.fail(c, err)
		return
	}
	var code any
	if q.CouponCode != nil {
		code = *q.CouponCode
	}
	c.JSON(http.StatusOK, gin.H{
		"items": q.Items, "subtotal": q.Subtotal, "coupon_code": code, "coupon_discount": q.Discount,
		"taxable_amount": q.Taxable, "gst_amount": q.GST, "total": q.Total, "address": addressJSON(q.Address),
	})
}

// ValidateCoupon godoc POST /api/v1/checkout/coupons/validate
func (h *CheckoutHandler) ValidateCoupon(c *gin.Context) {
	var req couponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, "coupon_code is required")
		return
	}
	v, err := h.svc.ValidateCoupon(c.Request.Context(), currentUser(c), req.CouponCode)
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"is_valid": v.Valid, "coupon_code": v.Code, "discount_type": v.Type,
		"discount_value": v.Value, "message": v.Message,
	})
}

// PlaceOrder godoc POST /api/v1/checkout/orders
func (h *CheckoutHandler) PlaceOrder(c *gin.Context) {
	var req checkoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, "address_id is required")
		return
	}
	v, err := h.svc.PlaceOrder(c.Request.Context(), currentUser(c), req.AddressID, req.CouponCode)
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, orderJSON(v))
}
