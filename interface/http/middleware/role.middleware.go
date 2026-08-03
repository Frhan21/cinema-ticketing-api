package middleware

import (
	"cinema-ticketing-api/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")

		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, response.ErrorResponse("Role not found"))
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, response.ErrorResponse("Invalid role type"))
			return
		}

		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, response.ErrorResponse("Forbidden: insufficient permissions"))
	}
}
