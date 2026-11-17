package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/models"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

// AdminAttributeHandler serves /api/v1/admin/attributes endpoints.
type AdminAttributeHandler struct {
	svc *services.AttributeAdminService
}

// NewAdminAttributeHandler returns an AdminAttributeHandler.
func NewAdminAttributeHandler(svc *services.AttributeAdminService) *AdminAttributeHandler {
	return &AdminAttributeHandler{svc: svc}
}

type attrDefRequest struct {
	Name         *string `json:"name"`
	IsFilterable *bool   `json:"is_filterable"`
	SortOrder    *int32  `json:"sort_order"`
	Code         string  `json:"code"`
}

type attrValueRequest struct {
	Value     *string `json:"value"`
	ValueSlug *string `json:"value_slug"`
	SortOrder *int32  `json:"sort_order"`
}

func defJSON(a *dbsqlc.ProductAttribute) gin.H {
	return gin.H{
		"id": a.ID, "code": a.Code, "name": a.Name, "is_filterable": a.IsFilterable,
		"sort_order": a.SortOrder, "created_at": a.CreatedAt, "updated_at": a.UpdatedAt,
	}
}

func valueModel(v *dbsqlc.ProductAttributeValue) models.AttributeValue {
	return models.AttributeValue{
		ID: v.ID, AttributeID: v.AttributeID, Value: v.Value, ValueSlug: v.Slug,
		SortOrder: v.SortOrder, CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt,
	}
}

func (h *AdminAttributeHandler) fail(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrNotFound):
		apiError(c, http.StatusNotFound, "attribute not found")
	case errors.Is(err, services.ErrConflict):
		apiError(c, http.StatusConflict, "conflict: duplicate code/slug, or value still assigned to products")
	case errors.Is(err, services.ErrValidation):
		apiError(c, http.StatusUnprocessableEntity, "validation failed")
	default:
		apiError(c, http.StatusInternalServerError, "internal error")
	}
}

func uuidParam(c *gin.Context, name string) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param(name))
	if err != nil {
		apiError(c, http.StatusNotFound, "attribute not found")
		return uuid.Nil, false
	}
	return id, true
}

// ListDefinitions godoc GET /api/v1/admin/attributes
func (h *AdminAttributeHandler) ListDefinitions(c *gin.Context) {
	page, err1 := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, err2 := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err1 != nil || err2 != nil || page < 1 || page > 1_000_000 || limit < 1 || limit > 100 {
		apiError(c, http.StatusBadRequest, "invalid pagination")
		return
	}
	rows, total, err := h.svc.ListDefinitions(c.Request.Context(), int32(page), int32(limit)) //nolint:gosec // bounded above
	if err != nil {
		h.fail(c, err)
		return
	}
	data := make([]gin.H, 0, len(rows))
	for i := range rows {
		data = append(data, defJSON(&rows[i]))
	}
	totalPages := (int(total) + limit - 1) / limit
	c.JSON(http.StatusOK, gin.H{"data": data, "meta": gin.H{"total": total, "page": page, "limit": limit, "total_pages": totalPages}})
}

// CreateDefinition godoc POST /api/v1/admin/attributes
func (h *AdminAttributeHandler) CreateDefinition(c *gin.Context) {
	var req attrDefRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	a, err := h.svc.CreateDefinition(c.Request.Context(), services.AttributeDefInput{
		Code: req.Code, Name: req.Name, IsFilterable: req.IsFilterable, SortOrder: req.SortOrder,
	})
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, defJSON(a))
}

// GetDefinition godoc GET /api/v1/admin/attributes/:attrDefId
func (h *AdminAttributeHandler) GetDefinition(c *gin.Context) {
	id, ok := uuidParam(c, "attrDefId")
	if !ok {
		return
	}
	a, err := h.svc.GetDefinition(c.Request.Context(), id)
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, defJSON(a))
}

// UpdateDefinition godoc PUT /api/v1/admin/attributes/:attrDefId
func (h *AdminAttributeHandler) UpdateDefinition(c *gin.Context) {
	id, ok := uuidParam(c, "attrDefId")
	if !ok {
		return
	}
	var req attrDefRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	a, err := h.svc.UpdateDefinition(c.Request.Context(), id, services.AttributeDefInput{
		Name: req.Name, IsFilterable: req.IsFilterable, SortOrder: req.SortOrder,
	})
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, defJSON(a))
}

// DeleteDefinition godoc DELETE /api/v1/admin/attributes/:attrDefId
func (h *AdminAttributeHandler) DeleteDefinition(c *gin.Context) {
	id, ok := uuidParam(c, "attrDefId")
	if !ok {
		return
	}
	if err := h.svc.DeleteDefinition(c.Request.Context(), id); err != nil {
		h.fail(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// ListValues godoc GET /api/v1/admin/attributes/:attrDefId/values
func (h *AdminAttributeHandler) ListValues(c *gin.Context) {
	id, ok := uuidParam(c, "attrDefId")
	if !ok {
		return
	}
	rows, err := h.svc.ListValues(c.Request.Context(), id)
	if err != nil {
		h.fail(c, err)
		return
	}
	data := make([]models.AttributeValue, 0, len(rows))
	for i := range rows {
		data = append(data, valueModel(&rows[i]))
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// CreateValue godoc POST /api/v1/admin/attributes/:attrDefId/values
func (h *AdminAttributeHandler) CreateValue(c *gin.Context) {
	id, ok := uuidParam(c, "attrDefId")
	if !ok {
		return
	}
	var req attrValueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	v, err := h.svc.CreateValue(c.Request.Context(), id, services.AttributeValueInput{Value: req.Value, ValueSlug: req.ValueSlug, SortOrder: req.SortOrder})
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, valueModel(v))
}

// UpdateValue godoc PUT /api/v1/admin/attributes/:attrDefId/values/:valueId
func (h *AdminAttributeHandler) UpdateValue(c *gin.Context) {
	id, ok := uuidParam(c, "attrDefId")
	if !ok {
		return
	}
	vid, ok := uuidParam(c, "valueId")
	if !ok {
		return
	}
	var req attrValueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	v, err := h.svc.UpdateValue(c.Request.Context(), id, vid, services.AttributeValueInput{Value: req.Value, ValueSlug: req.ValueSlug, SortOrder: req.SortOrder})
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, valueModel(v))
}

// DeleteValue godoc DELETE /api/v1/admin/attributes/:attrDefId/values/:valueId
func (h *AdminAttributeHandler) DeleteValue(c *gin.Context) {
	id, ok := uuidParam(c, "attrDefId")
	if !ok {
		return
	}
	vid, ok := uuidParam(c, "valueId")
	if !ok {
		return
	}
	if err := h.svc.DeleteValue(c.Request.Context(), id, vid); err != nil {
		h.fail(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
