package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

// AdminOrderHandler serves /api/v1/admin/orders endpoints.
type AdminOrderHandler struct {
	svc *services.AdminOrderService
}

// NewAdminOrderHandler returns an AdminOrderHandler.
func NewAdminOrderHandler(svc *services.AdminOrderService) *AdminOrderHandler {
	return &AdminOrderHandler{svc: svc}
}

func adminOrderJSON(v *services.AdminOrderView) gin.H {
	j := orderJSON(&v.OrderView)
	j["customer_name"] = v.CustomerName
	j["customer_email"] = v.CustomerEmail
	return j
}

func (h *AdminOrderHandler) fail(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrNotFound):
		apiError(c, http.StatusNotFound, "order not found")
	case errors.Is(err, services.ErrInvalidStatus):
		apiError(c, http.StatusBadRequest, "invalid status filter")
	case errors.Is(err, services.ErrInvalidTransition):
		apiError(c, http.StatusUnprocessableEntity, "order cannot move to that status")
	default:
		apiError(c, http.StatusInternalServerError, "internal error")
	}
}

// List godoc GET /api/v1/admin/orders
func (h *AdminOrderHandler) List(c *gin.Context) {
	page, err1 := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, err2 := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err1 != nil || err2 != nil || page < 1 || page > 1_000_000 || limit < 1 || limit > 100 {
		apiError(c, http.StatusBadRequest, "invalid pagination")
		return
	}
	res, err := h.svc.List(c.Request.Context(), c.Query("status"), page, limit)
	if err != nil {
		h.fail(c, err)
		return
	}
	data := make([]gin.H, 0, len(res.Orders))
	for i := range res.Orders {
		data = append(data, adminOrderJSON(&res.Orders[i]))
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "meta": gin.H{
		"total": res.Total, "page": page, "limit": limit, "total_pages": (int(res.Total) + limit - 1) / limit,
	}})
}

// Get godoc GET /api/v1/admin/orders/:id
func (h *AdminOrderHandler) Get(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	v, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, adminOrderJSON(v))
}

// UpdateStatus godoc PUT /api/v1/admin/orders/:id/status
func (h *AdminOrderHandler) UpdateStatus(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	v, err := h.svc.UpdateStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, adminOrderJSON(v))
}
