package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

// AttributeHandler serves the public attribute (filter option) endpoints.
type AttributeHandler struct {
	svc *services.AttributeService
}

// NewAttributeHandler returns an AttributeHandler.
func NewAttributeHandler(svc *services.AttributeService) *AttributeHandler {
	return &AttributeHandler{svc: svc}
}

// List godoc GET /api/v1/attributes
func (h *AttributeHandler) List(c *gin.Context) {
	data, err := h.svc.ListFilterable(c.Request.Context())
	if err != nil {
		apiError(c, http.StatusInternalServerError, "internal error")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// ListValues godoc GET /api/v1/attributes/:attrDefId/values
func (h *AttributeHandler) ListValues(c *gin.Context) {
	id, err := uuid.Parse(c.Param("attrDefId"))
	if err != nil {
		apiError(c, http.StatusNotFound, "attribute not found")
		return
	}
	data, err := h.svc.ListValues(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			apiError(c, http.StatusNotFound, "attribute not found")
			return
		}
		apiError(c, http.StatusInternalServerError, "internal error")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}
