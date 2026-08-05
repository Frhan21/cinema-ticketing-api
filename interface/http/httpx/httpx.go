package httpx

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetUserIDFromContext mengambil dan mem-parse user_id dari Gin context
// yang sebelumnya diset oleh AuthMiddleware.
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	raw, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}

	str, ok := raw.(string)
	if !ok {
		return uuid.Nil, false
	}

	parsed, err := uuid.Parse(str)
	if err != nil {
		return uuid.Nil, false
	}

	return parsed, true
}

func GetRoleFromContext(c *gin.Context) (string, bool) {
	role, exists := c.Get("role")
	if !exists {
		return "", false
	}
	value, ok := role.(string)
	return value, ok
}
