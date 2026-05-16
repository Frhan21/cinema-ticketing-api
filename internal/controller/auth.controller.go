package controller

import (
	"cinema-ticketing-api/internal/request"
	"cinema-ticketing-api/internal/response"
	"cinema-ticketing-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (ac *AuthController) Register(c *gin.Context) {
	var req request.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	res, err := ac.authService.Register(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("You've been Registered", res))
}

func (ac *AuthController) Login(c *gin.Context) {
	var req request.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	res, err := ac.authService.Login(req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("You've been Logged In", res))

}
