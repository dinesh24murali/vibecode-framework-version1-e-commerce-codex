package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

// PaymentHandler serves payment initiation and the Razorpay webhook.
type PaymentHandler struct {
	svc *services.PaymentService
}

// NewPaymentHandler returns a PaymentHandler.
func NewPaymentHandler(svc *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

type initiatePaymentRequest struct {
	OrderID uuid.UUID `json:"order_id" binding:"required"`
}

// Initiate godoc POST /api/v1/payments/initiate
func (h *PaymentHandler) Initiate(c *gin.Context) {
	var req initiatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, "order_id is required")
		return
	}
	out, err := h.svc.Initiate(c.Request.Context(), currentUser(c), req.OrderID, c.GetHeader("Idempotency-Key"))
	switch {
	case err == nil:
		c.JSON(http.StatusOK, out)
	case errors.Is(err, services.ErrMissingIdemKey):
		apiError(c, http.StatusBadRequest, "Idempotency-Key header is required (8-128 characters)")
	case errors.Is(err, services.ErrNotFound):
		apiError(c, http.StatusNotFound, "order not found")
	case errors.Is(err, services.ErrOrderNotPayable), errors.Is(err, services.ErrIdempotencyInFlight):
		apiError(c, http.StatusConflict, err.Error())
	case errors.Is(err, services.ErrPaymentsDisabled):
		apiError(c, http.StatusServiceUnavailable, "payments are not configured")
	default:
		apiError(c, http.StatusInternalServerError, "internal error")
	}
}

// Webhook godoc POST /api/v1/payments/webhooks/razorpay
func (h *PaymentHandler) Webhook(c *gin.Context) {
	body, err := io.ReadAll(io.LimitReader(c.Request.Body, 1<<20))
	if err != nil {
		apiError(c, http.StatusBadRequest, "unreadable body")
		return
	}
	err = h.svc.HandleWebhook(c.Request.Context(), body, c.GetHeader("X-Razorpay-Signature"))
	switch {
	case err == nil:
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	case errors.Is(err, services.ErrInvalidSignature):
		apiError(c, http.StatusBadRequest, "invalid signature")
	default:
		apiError(c, http.StatusInternalServerError, "internal error")
	}
}
