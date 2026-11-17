package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

// AdminCategoryHandler serves /api/v1/admin/categories endpoints.
type AdminCategoryHandler struct {
	svc *services.AdminCategoryService
}

// NewAdminCategoryHandler returns an AdminCategoryHandler.
func NewAdminCategoryHandler(svc *services.AdminCategoryService) *AdminCategoryHandler {
	return &AdminCategoryHandler{svc: svc}
}

type adminCategoryRequest struct {
	Name     *string `json:"name"`
	Slug     *string `json:"slug"`
	ImageURL *string `json:"image_url"`
}

func (r adminCategoryRequest) input() services.AdminCategoryInput {
	return services.AdminCategoryInput{Name: r.Name, Slug: r.Slug, ImageURL: r.ImageURL}
}

func (h *AdminCategoryHandler) fail(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrNotFound):
		apiError(c, http.StatusNotFound, "category not found")
	case errors.Is(err, services.ErrConflict):
		apiError(c, http.StatusConflict, "conflict: slug already in use, or the category still has products")
	case errors.Is(err, services.ErrValidation):
		apiError(c, http.StatusUnprocessableEntity, "validation failed")
	default:
		apiError(c, http.StatusInternalServerError, "internal error")
	}
}

// List godoc GET /api/v1/admin/categories
func (h *AdminCategoryHandler) List(c *gin.Context) {
	page, err1 := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, err2 := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err1 != nil || err2 != nil || page < 1 || page > 1_000_000 || limit < 1 || limit > 100 {
		apiError(c, http.StatusBadRequest, "invalid pagination")
		return
	}
	data, total, err := h.svc.List(c.Request.Context(), page, limit)
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "meta": gin.H{
		"total": total, "page": page, "limit": limit, "total_pages": (int(total) + limit - 1) / limit,
	}})
}

// Create godoc POST /api/v1/admin/categories
func (h *AdminCategoryHandler) Create(c *gin.Context) {
	var req adminCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	cat, err := h.svc.Create(c.Request.Context(), req.input())
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, cat)
}

// Get godoc GET /api/v1/admin/categories/:id
func (h *AdminCategoryHandler) Get(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	cat, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, cat)
}

// Update godoc PUT /api/v1/admin/categories/:id
func (h *AdminCategoryHandler) Update(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	var req adminCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	cat, err := h.svc.Update(c.Request.Context(), id, req.input())
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, cat)
}

// Delete godoc DELETE /api/v1/admin/categories/:id
func (h *AdminCategoryHandler) Delete(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "category deleted"})
}
