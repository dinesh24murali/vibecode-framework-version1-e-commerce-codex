package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// ContextUserID is the gin context key for the authenticated user's UUID.
	ContextUserID = "userID"
	// ContextUserRole is the gin context key for the authenticated user's role.
	ContextUserRole = "userRole"
)

// Middleware returns a Gin handler that validates the Bearer access token.
// On success it sets ContextUserID and ContextUserRole in the request context.
// On failure it aborts with 401 Unauthorized.
func Middleware(jwtSvc *JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header format"})
			return
		}

		claims, err := jwtSvc.VerifyAccessToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token subject"})
			return
		}

		c.Set(ContextUserID, userID)
		c.Set(ContextUserRole, claims.Role)
		c.Next()
	}
}
