package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

// AdminCouponHandler serves /api/v1/admin/coupons endpoints.
type AdminCouponHandler struct {
	svc *services.AdminCouponService
}

// NewAdminCouponHandler returns an AdminCouponHandler.
func NewAdminCouponHandler(svc *services.AdminCouponService) *AdminCouponHandler {
	return &AdminCouponHandler{svc: svc}
}

// optTime decodes an absent / null / RFC3339 value so that "clear the expiry" is distinguishable.
type optTime services.NullableTime

func (o *optTime) UnmarshalJSON(b []byte) error {
	o.Set = true
	if bytes.Equal(bytes.TrimSpace(b), []byte("null")) {
		return nil
	}
	t, err := time.Parse(time.RFC3339, string(bytes.Trim(b, `"`)))
	if err != nil {
		return err
	}
	o.Time, o.Valid = t, true
	return nil
}

type adminCouponRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Code        *string `json:"code"`
	IsEnable    *bool   `json:"is_enable"`
	Type        *string `json:"type"`
	Flat        *int32  `json:"flat"`
	Percent     *int32  `json:"percent"`
	Expiry      optTime `json:"coupon_expiry"`
}

func (r adminCouponRequest) input() services.AdminCouponInput {
	return services.AdminCouponInput{
		Name: r.Name, Description: r.Description, Code: r.Code, IsEnable: r.IsEnable, Type: r.Type,
		Flat: r.Flat, Percent: r.Percent, Expiry: services.NullableTime(r.Expiry),
	}
}

func couponJSON(c *dbsqlc.Coupon) gin.H {
	var desc, flat, percent, expiry any
	if c.Description.Valid {
		desc = c.Description.String
	}
	if c.Flat.Valid {
		flat = c.Flat.Int32
	}
	if c.Percent.Valid {
		percent = c.Percent.Int32
	}
	if c.CouponExpiry.Valid {
		expiry = c.CouponExpiry.Time
	}
	return gin.H{
		"id": c.ID, "name": c.Name, "description": desc, "code": c.Code, "is_enable": c.IsEnable,
		"type": c.Type, "flat": flat, "percent": percent, "coupon_expiry": expiry,
		"created_at": c.CreatedAt, "updated_at": c.UpdatedAt,
	}
}

func (h *AdminCouponHandler) fail(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrNotFound):
		apiError(c, http.StatusNotFound, "coupon not found")
	case errors.Is(err, services.ErrConflict):
		apiError(c, http.StatusConflict, "conflict: code already in use, or the coupon has been redeemed (disable it instead)")
	case errors.Is(err, services.ErrValidation):
		apiError(c, http.StatusUnprocessableEntity, "validation failed")
	default:
		apiError(c, http.StatusInternalServerError, "internal error")
	}
}

// List godoc GET /api/v1/admin/coupons
func (h *AdminCouponHandler) List(c *gin.Context) {
	page, err1 := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, err2 := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err1 != nil || err2 != nil || page < 1 || page > 1_000_000 || limit < 1 || limit > 100 {
		apiError(c, http.StatusBadRequest, "invalid pagination")
		return
	}
	rows, total, err := h.svc.List(c.Request.Context(), page, limit)
	if err != nil {
		h.fail(c, err)
		return
	}
	data := make([]gin.H, 0, len(rows))
	for i := range rows {
		data = append(data, couponJSON(&rows[i]))
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "meta": gin.H{
		"total": total, "page": page, "limit": limit, "total_pages": (int(total) + limit - 1) / limit,
	}})
}

// Create godoc POST /api/v1/admin/coupons
func (h *AdminCouponHandler) Create(c *gin.Context) {
	var req adminCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	cp, err := h.svc.Create(c.Request.Context(), req.input())
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, couponJSON(cp))
}

// Get godoc GET /api/v1/admin/coupons/:id
func (h *AdminCouponHandler) Get(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	cp, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, couponJSON(cp))
}

// Update godoc PUT /api/v1/admin/coupons/:id
func (h *AdminCouponHandler) Update(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	var req adminCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	cp, err := h.svc.Update(c.Request.Context(), id, req.input())
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, couponJSON(cp))
}

// Delete godoc DELETE /api/v1/admin/coupons/:id
func (h *AdminCouponHandler) Delete(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "coupon deleted"})
}
