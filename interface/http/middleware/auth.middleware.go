package middleware

import (
	"cinema-ticketing-api/config"
	"cinema-ticketing-api/pkg/jwt"
	"cinema-ticketing-api/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtCfg config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse("Unauthorized"))
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse("Invalid authorization format"))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwt.ValidateToken(tokenString, jwtCfg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse("Invalid or expired token"))
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
