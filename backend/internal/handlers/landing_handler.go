package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

// LandingHandler serves the public landing page payload.
type LandingHandler struct {
	svc *services.LandingService
}

// NewLandingHandler returns a LandingHandler.
func NewLandingHandler(svc *services.LandingService) *LandingHandler {
	return &LandingHandler{svc: svc}
}

// Get godoc GET /api/v1/landing
func (h *LandingHandler) Get(c *gin.Context) {
	res, err := h.svc.Landing(c.Request.Context())
	if err != nil {
		apiError(c, http.StatusInternalServerError, "internal error")
		return
	}
	c.JSON(http.StatusOK, res)
}
