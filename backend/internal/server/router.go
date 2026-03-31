package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/handlers"
)

// newRouter builds and returns a configured *gin.Engine.
// It applies global middleware (CORS, recovery, structured logging) and
// registers all routes. Feature route groups are added here as tasks
// implement them. The /health endpoint is registered at root level so
// infrastructure health checks do not traverse auth middleware.
func newRouter(allowedOrigins []string) *gin.Engine {
	r := gin.New()

	// Recovery middleware: turns panics into 500 responses and logs the stack.
	r.Use(gin.Recovery())

	// Structured request logger — uses gin's built-in JSON-compatible logger.
	// Replace with zerolog middleware in a follow-up task if desired.
	r.Use(gin.Logger())

	// CORS — exact origin matching, no wildcards (architecture constraint).
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Internal health check — outside /api/v1 and auth middleware intentionally.
	r.GET("/health", handlers.Health)

	// API v1 group — feature routes are added here in subsequent tasks.
	// v1 := r.Group("/api/v1")

	return r
}
