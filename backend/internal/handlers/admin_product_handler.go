package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

// AdminProductHandler serves /api/v1/admin/products endpoints.
type AdminProductHandler struct {
	svc *services.AdminProductService
}

// NewAdminProductHandler returns an AdminProductHandler.
func NewAdminProductHandler(svc *services.AdminProductService) *AdminProductHandler {
	return &AdminProductHandler{svc: svc}
}

type adminImageRequest struct {
	AltText   *string `json:"alt_text"`
	URL       string  `json:"url"`
	SortOrder int32   `json:"sort_order"`
}

type adminProductRequest struct {
	Name            *string              `json:"name"`
	Slug            *string              `json:"slug"`
	Description     *string              `json:"description"`
	Price           *float64             `json:"price"`
	DiscountPercent *float64             `json:"discount_percent"`
	Stock           *int32               `json:"stock"`
	IsActive        *bool                `json:"is_active"`
	CategoryIDs     *[]uuid.UUID         `json:"category_ids"`
	ImageURLs       *[]adminImageRequest `json:"image_urls"`
}

func (r adminProductRequest) input() services.AdminProductInput {
	in := services.AdminProductInput{
		Name: r.Name, Slug: r.Slug, Description: r.Description, Price: r.Price,
		DiscountPercent: r.DiscountPercent, Stock: r.Stock, IsActive: r.IsActive, CategoryIDs: r.CategoryIDs,
	}
	if r.ImageURLs != nil {
		imgs := make([]services.AdminImageInput, 0, len(*r.ImageURLs))
		for _, i := range *r.ImageURLs {
			imgs = append(imgs, services.AdminImageInput{URL: i.URL, AltText: i.AltText, SortOrder: i.SortOrder})
		}
		in.Images = &imgs
	}
	return in
}

func (h *AdminProductHandler) fail(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrNotFound):
		apiError(c, http.StatusNotFound, "not found")
	case errors.Is(err, services.ErrConflict):
		apiError(c, http.StatusConflict, "conflict: slug already in use")
	case errors.Is(err, services.ErrValidation):
		apiError(c, http.StatusUnprocessableEntity, "validation failed")
	default:
		apiError(c, http.StatusInternalServerError, "internal error")
	}
}

// List godoc GET /api/v1/admin/products
func (h *AdminProductHandler) List(c *gin.Context) {
	page, err1 := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, err2 := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err1 != nil || err2 != nil || page < 1 || page > 1_000_000 || limit < 1 || limit > 100 {
		apiError(c, http.StatusBadRequest, "invalid pagination")
		return
	}
	data, total, err := h.svc.List(c.Request.Context(), c.Query("q"), page, limit)
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "meta": gin.H{
		"total": total, "page": page, "limit": limit, "total_pages": (int(total) + limit - 1) / limit,
	}})
}

// Create godoc POST /api/v1/admin/products
func (h *AdminProductHandler) Create(c *gin.Context) {
	var req adminProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	p, err := h.svc.Create(c.Request.Context(), req.input())
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, p)
}

// Get godoc GET /api/v1/admin/products/:id
func (h *AdminProductHandler) Get(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	p, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, p)
}

// Update godoc PUT /api/v1/admin/products/:id
func (h *AdminProductHandler) Update(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	var req adminProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	p, err := h.svc.Update(c.Request.Context(), id, req.input())
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, p)
}

// Delete godoc DELETE /api/v1/admin/products/:id (soft delete)
func (h *AdminProductHandler) Delete(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}

// ListAttributes godoc GET /api/v1/admin/products/:id/attributes
func (h *AdminProductHandler) ListAttributes(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	data, err := h.svc.ListAttributes(c.Request.Context(), id)
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// AssignAttribute godoc POST /api/v1/admin/products/:id/attributes
func (h *AdminProductHandler) AssignAttribute(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	var req struct {
		AttributeValueID uuid.UUID `json:"attribute_value_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	a, err := h.svc.AssignAttribute(c.Request.Context(), id, req.AttributeValueID)
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, a)
}

// RemoveAttribute godoc DELETE /api/v1/admin/products/:id/attributes/:attrId
func (h *AdminProductHandler) RemoveAttribute(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	valueID, err := uuid.Parse(c.Param("attrId"))
	if err != nil {
		apiError(c, http.StatusNotFound, "not found")
		return
	}
	if err := h.svc.RemoveAttribute(c.Request.Context(), id, valueID); err != nil {
		h.fail(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
