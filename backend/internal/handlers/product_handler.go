package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

// ProductHandler serves the public product endpoints.
type ProductHandler struct {
	svc *services.ProductService
}

// NewProductHandler returns a ProductHandler.
func NewProductHandler(svc *services.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

// queryArray reads a repeated query param, accepting both name[] and name.
func queryArray(c *gin.Context, name string) []string {
	vals := append(c.QueryArray(name+"[]"), c.QueryArray(name)...)
	out := make([]string, 0, len(vals))
	for _, v := range vals {
		if v = strings.TrimSpace(v); v != "" {
			out = append(out, v)
		}
	}
	return out
}

// List godoc GET /api/v1/products
func (h *ProductHandler) List(c *gin.Context) {
	page, err1 := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, err2 := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err1 != nil || err2 != nil || page < 1 || limit < 1 || limit > 100 {
		apiError(c, http.StatusBadRequest, "invalid pagination parameters")
		return
	}
	var ids []uuid.UUID
	for _, raw := range queryArray(c, "attr_value_id") {
		id, err := uuid.Parse(raw)
		if err != nil {
			apiError(c, http.StatusBadRequest, "attr_value_id must be a UUID")
			return
		}
		ids = append(ids, id)
	}
	res, err := h.svc.ListProducts(c.Request.Context(), services.ProductQuery{
		Query: strings.TrimSpace(c.Query("q")), Sort: c.DefaultQuery("sort_by", "newest"),
		CategorySlugs: queryArray(c, "category"), ValueIDs: ids, Page: page, Limit: limit,
	})
	if err != nil {
		if errors.Is(err, services.ErrInvalidSort) {
			apiError(c, http.StatusBadRequest, err.Error())
			return
		}
		apiError(c, http.StatusInternalServerError, "internal error")
		return
	}
	c.JSON(http.StatusOK, res)
}

// Get godoc GET /api/v1/products/:id
func (h *ProductHandler) Get(c *gin.Context) {
	id, ok := pathID(c)
	if !ok {
		return
	}
	p, err := h.svc.GetProduct(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			apiError(c, http.StatusNotFound, "product not found")
			return
		}
		apiError(c, http.StatusInternalServerError, "internal error")
		return
	}
	c.JSON(http.StatusOK, p)
}
