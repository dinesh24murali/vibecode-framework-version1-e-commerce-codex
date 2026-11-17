package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

// CategoryHandler serves the public category endpoints.
type CategoryHandler struct {
	svc *services.ProductService
}

// NewCategoryHandler returns a CategoryHandler.
func NewCategoryHandler(svc *services.ProductService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

// List godoc GET /api/v1/categories
func (h *CategoryHandler) List(c *gin.Context) {
	data, err := h.svc.ListCategories(c.Request.Context())
	if err != nil {
		apiError(c, http.StatusInternalServerError, "internal error")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// Get godoc GET /api/v1/categories/:id
func (h *CategoryHandler) Get(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	cat, err := h.svc.GetCategory(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			apiError(c, http.StatusNotFound, "category not found")
			return
		}
		apiError(c, http.StatusInternalServerError, "internal error")
		return
	}
	c.JSON(http.StatusOK, cat)
}
