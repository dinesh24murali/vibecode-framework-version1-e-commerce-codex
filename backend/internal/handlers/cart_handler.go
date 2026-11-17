package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

// CartHandler serves the authenticated cart endpoints.
type CartHandler struct {
	svc *services.CartService
}

// NewCartHandler returns a CartHandler.
func NewCartHandler(svc *services.CartService) *CartHandler {
	return &CartHandler{svc: svc}
}

type addCartRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int32     `json:"quantity"`
}

type updateCartRequest struct {
	Quantity int32 `json:"quantity"`
}

func (h *CartHandler) fail(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrInvalidQuantity):
		apiError(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, services.ErrNotFound):
		apiError(c, http.StatusNotFound, "not found")
	default:
		apiError(c, http.StatusInternalServerError, "internal error")
	}
}

// Get godoc GET /api/v1/cart
func (h *CartHandler) Get(c *gin.Context) {
	cart, err := h.svc.Cart(c.Request.Context(), currentUser(c))
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, cart)
}

// Clear godoc DELETE /api/v1/cart
func (h *CartHandler) Clear(c *gin.Context) {
	if err := h.svc.Clear(c.Request.Context(), currentUser(c)); err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cart cleared."})
}

// AddItem godoc POST /api/v1/cart/items
func (h *CartHandler) AddItem(c *gin.Context) {
	var req addCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, "product_id and quantity are required")
		return
	}
	cart, err := h.svc.Add(c.Request.Context(), currentUser(c), req.ProductID, req.Quantity)
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, cart)
}

// UpdateItem godoc PUT /api/v1/cart/items/:id
func (h *CartHandler) UpdateItem(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	var req updateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, "quantity is required")
		return
	}
	cart, err := h.svc.SetQuantity(c.Request.Context(), currentUser(c), id, req.Quantity)
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, cart)
}

// RemoveItem godoc DELETE /api/v1/cart/items/:id
func (h *CartHandler) RemoveItem(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	cart, err := h.svc.Remove(c.Request.Context(), currentUser(c), id)
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, cart)
}
