package controller

import (
	"cinema-ticketing-api/internal/response"
	"cinema-ticketing-api/pkg/mailer"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TestSendEmail(c *gin.Context) {
	mail := mailer.NewMailer()
	err := mail.SendEmail("frhn.r3@gmail.com", "Test Email Cinema Ticketing", "<h1>Hello</h1><p>Email berhasil dikirim dari Go.</p>")
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Email berhasil dikirim", nil))
}
