package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleAdmin is the JWT role claim value for administrators.
const RoleAdmin = "admin"

// AdminMiddleware authenticates the Bearer access token (401 on failure) and
// then requires the `role: admin` claim (403 for any other role).
func AdminMiddleware(jwtSvc *JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !authenticate(c, jwtSvc) {
			return
		}
		if role, _ := c.Get(ContextUserRole); role != RoleAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": http.StatusForbidden, "message": "forbidden"})
			return
		}
		c.Next()
	}
}
