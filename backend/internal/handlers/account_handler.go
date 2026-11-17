package handlers

import (
	"database/sql"
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/auth"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

// AccountHandler serves /account/* and /orders/* endpoints for the signed-in customer.
type AccountHandler struct {
	svc *services.AccountService
}

// NewAccountHandler returns an AccountHandler.
func NewAccountHandler(svc *services.AccountService) *AccountHandler {
	return &AccountHandler{svc: svc}
}

type updateProfileRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type changePasswordRequest struct {
	CurrentPassword    string `json:"current_password" binding:"required"`
	NewPassword        string `json:"new_password" binding:"required,min=8"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required"`
}

type requestEmailChangeRequest struct {
	NewEmail string `json:"new_email" binding:"required,email"`
}

type confirmEmailChangeRequest struct {
	NewEmail string `json:"new_email" binding:"required,email"`
	OTP      string `json:"otp" binding:"required"`
}

type addressRequest struct {
	AddressLine1 string  `json:"address_line1" binding:"required"`
	AddressLine2 *string `json:"address_line2"`
	Landmark     *string `json:"landmark"`
	City         string  `json:"city" binding:"required"`
	State        string  `json:"state" binding:"required"`
	Pincode      string  `json:"pincode" binding:"required"`
	CompanyName  *string `json:"company_name"`
	GSTNumber    *string `json:"gst_number"`
	Phone        string  `json:"phone" binding:"required"`
	IsDefault    bool    `json:"is_default"`
}

func (r addressRequest) input() services.AddressInput {
	return services.AddressInput{
		AddressLine1: r.AddressLine1, AddressLine2: r.AddressLine2, Landmark: r.Landmark,
		City: r.City, State: r.State, Pincode: r.Pincode, CompanyName: r.CompanyName,
		GSTNumber: r.GSTNumber, Phone: r.Phone, IsDefault: r.IsDefault,
	}
}

// apiError writes the spec ErrorResponse shape: {code, message}.
func apiError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"code": status, "message": message})
}

func currentUser(c *gin.Context) uuid.UUID {
	id, _ := c.Get(auth.ContextUserID)
	uid, _ := id.(uuid.UUID)
	return uid
}

func pathID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiError(c, http.StatusBadRequest, "invalid id")
		return uuid.Nil, false
	}
	return id, true
}

// fail maps service errors to HTTP responses.
func (h *AccountHandler) fail(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrNotFound):
		apiError(c, http.StatusNotFound, "not found")
	case errors.Is(err, services.ErrForbidden):
		apiError(c, http.StatusForbidden, "forbidden")
	case errors.Is(err, services.ErrWrongPassword), errors.Is(err, services.ErrPasswordMismatch),
		errors.Is(err, services.ErrInvalidOTP):
		apiError(c, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, services.ErrEmailAlreadyExists):
		apiError(c, http.StatusConflict, err.Error())
	case errors.Is(err, services.ErrOrderNotCancel), errors.Is(err, services.ErrInvalidStatus):
		apiError(c, http.StatusBadRequest, err.Error())
	default:
		apiError(c, http.StatusInternalServerError, "internal error")
	}
}

func profileJSON(u *dbsqlc.User) gin.H {
	first, last := splitFullName(u.FullName)
	return gin.H{
		"id": u.ID, "email": u.Email, "first_name": first, "last_name": last,
		"role": u.Role, "created_at": u.CreatedAt, "updated_at": u.UpdatedAt,
	}
}

func nullable(s sql.NullString) any {
	if s.Valid {
		return s.String
	}
	return nil
}

func addressJSON(a *dbsqlc.Address) gin.H {
	return gin.H{
		"id": a.ID, "address_line1": a.Line1, "address_line2": nullable(a.Line2),
		"landmark": nullable(a.Landmark), "city": a.City, "state": a.State,
		"pincode": a.PostalCode, "company_name": nullable(a.CompanyName),
		"gst_number": nullable(a.GstNumber), "phone": a.PhoneE164.String,
		"is_default": a.IsDefault, "created_at": a.CreatedAt, "updated_at": a.UpdatedAt,
	}
}

func money(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func orderJSON(v *services.OrderView) gin.H {
	items := make([]gin.H, 0, len(v.Items))
	for _, it := range v.Items {
		var pid any
		if it.ProductID.Valid {
			pid = it.ProductID.UUID
		}
		items = append(items, gin.H{
			"id": it.ID, "product_id": pid, "product_name": it.ProductTitle,
			"product_slug": it.ProductSlug, "cover_image_url": nil, "quantity": it.Quantity,
			"unit_price": money(it.UnitPriceInr), "line_total": money(it.LineTotalInr),
		})
	}
	addr := gin.H{}
	if a := v.Address; a != nil {
		addr = gin.H{
			"id": a.ID, "address_line1": a.Line1, "address_line2": nullable(a.Line2),
			"landmark": nil, "city": a.City, "state": a.State, "pincode": a.PostalCode,
			"company_name": nil, "gst_number": nil, "phone": a.PhoneE164.String,
			"is_default": false, "created_at": a.CreatedAt, "updated_at": a.UpdatedAt,
		}
	}
	o := *v.Order
	return gin.H{
		"id": o.ID, "order_number": o.OrderNumber, "status": v.Status, "items": items,
		"shipping_address": addr, "subtotal": money(o.SubtotalInr), "coupon_code": nullable(o.CouponCode),
		"coupon_discount": money(o.CouponDiscountInr), "taxable_amount": money(o.SubtotalInr) - money(o.CouponDiscountInr), "gst_amount": money(o.TaxInr),
		"total": money(o.TotalInr), "created_at": o.CreatedAt, "updated_at": o.UpdatedAt,
	}
}

// ── Profile ──────────────────────────────────────────────────────────────────

// GetProfile godoc GET /api/v1/account/profile
func (h *AccountHandler) GetProfile(c *gin.Context) {
	u, err := h.svc.GetProfile(c.Request.Context(), currentUser(c))
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, profileJSON(u))
}

// UpdateProfile godoc PUT /api/v1/account/profile
func (h *AccountHandler) UpdateProfile(c *gin.Context) {
	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, err.Error())
		return
	}
	u, err := h.svc.UpdateProfile(c.Request.Context(), currentUser(c), req.FirstName, req.LastName)
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, profileJSON(u))
}

// ChangePassword godoc PUT /api/v1/account/profile/password
func (h *AccountHandler) ChangePassword(c *gin.Context) {
	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.svc.ChangePassword(c.Request.Context(), currentUser(c), req.CurrentPassword, req.NewPassword, req.ConfirmNewPassword); err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password updated."})
}

// RequestEmailChange godoc POST /api/v1/account/profile/email-change/request
func (h *AccountHandler) RequestEmailChange(c *gin.Context) {
	var req requestEmailChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, err.Error())
		return
	}
	h.svc.RequestEmailChange(c.Request.Context(), req.NewEmail)
	c.JSON(http.StatusOK, gin.H{"message": "If the address is available, a verification code has been sent."})
}

// ConfirmEmailChange godoc POST /api/v1/account/profile/email-change/confirm
func (h *AccountHandler) ConfirmEmailChange(c *gin.Context) {
	var req confirmEmailChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, err.Error())
		return
	}
	u, err := h.svc.ConfirmEmailChange(c.Request.Context(), currentUser(c), req.NewEmail, req.OTP)
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, profileJSON(u))
}

// ── Addresses ────────────────────────────────────────────────────────────────

// ListAddresses godoc GET /api/v1/account/addresses
func (h *AccountHandler) ListAddresses(c *gin.Context) {
	rows, err := h.svc.ListAddresses(c.Request.Context(), currentUser(c))
	if err != nil {
		h.fail(c, err)
		return
	}
	data := make([]gin.H, 0, len(rows))
	for i := range rows {
		data = append(data, addressJSON(&rows[i]))
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// CreateAddress godoc POST /api/v1/account/addresses
func (h *AccountHandler) CreateAddress(c *gin.Context) {
	var req addressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, err.Error())
		return
	}
	a, err := h.svc.CreateAddress(c.Request.Context(), currentUser(c), req.input())
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, addressJSON(a))
}

// UpdateAddress godoc PUT /api/v1/account/addresses/:id
func (h *AccountHandler) UpdateAddress(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	var req addressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, err.Error())
		return
	}
	a, err := h.svc.UpdateAddress(c.Request.Context(), currentUser(c), id, req.input())
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, addressJSON(a))
}

// DeleteAddress godoc DELETE /api/v1/account/addresses/:id
func (h *AccountHandler) DeleteAddress(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteAddress(c.Request.Context(), currentUser(c), id); err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Address deleted."})
}

// ── Orders ───────────────────────────────────────────────────────────────────

// ListOrders godoc GET /api/v1/orders
func (h *AccountHandler) ListOrders(c *gin.Context) {
	page, err1 := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, err2 := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err1 != nil || err2 != nil || page < 1 || limit < 1 || limit > 100 {
		apiError(c, http.StatusBadRequest, "invalid pagination parameters")
		return
	}
	res, err := h.svc.ListOrders(c.Request.Context(), currentUser(c), c.Query("status"), page, limit)
	if err != nil {
		h.fail(c, err)
		return
	}
	data := make([]gin.H, 0, len(res.Orders))
	for i := range res.Orders {
		data = append(data, orderJSON(&res.Orders[i]))
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "meta": gin.H{
		"total": res.Total, "page": page, "limit": limit,
		"total_pages": int(math.Ceil(float64(res.Total) / float64(limit))),
	}})
}

// GetOrder godoc GET /api/v1/orders/:id
func (h *AccountHandler) GetOrder(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	v, err := h.svc.GetOrder(c.Request.Context(), currentUser(c), id)
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, orderJSON(v))
}

// CancelOrder godoc PUT /api/v1/orders/:id/cancel
func (h *AccountHandler) CancelOrder(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	v, err := h.svc.CancelOrder(c.Request.Context(), currentUser(c), id)
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, orderJSON(v))
}
