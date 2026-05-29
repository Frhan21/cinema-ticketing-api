package middleware

import (
	"cinema-ticketing-api/internal/response"
	"cinema-ticketing-api/pkg/apperror"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if appErr, ok := err.(*apperror.AppError); ok {
				c.JSON(appErr.Code, response.ErrorResponse(appErr.Message))
				return
			}

			// Default to 500 if the error is not an AppError
			c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		}
	}
}
