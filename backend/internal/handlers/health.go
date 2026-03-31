// Package handlers contains Gin HTTP handlers.
// Handlers are responsible only for parsing/validating input and delegating
// to service layer functions. They never call the database directly.
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthResponse is the JSON body returned by the health endpoint.
type HealthResponse struct {
	Status string `json:"status"`
}

// Health handles GET /health and returns 200 OK with {"status":"ok"}.
// This endpoint is intentionally outside the /api/v1 prefix so that
// infrastructure health checks do not require auth tokens.
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{Status: "ok"})
}
