package controller

import (
	"cinema-ticketing-api/internal/request"
	"cinema-ticketing-api/internal/response"
	"cinema-ticketing-api/internal/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (uc *UserController) GetProfile(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Unauthorized"))
		return
	}

	id, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Invalid user session"))
		return
	}

	user, err := uc.userService.GetProfile(id)
	if err != nil {
		c.JSON(userStatusCode(err), response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("User profile fetched successfully", user))
}

func (uc *UserController) UpdateProfile(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Unauthorized"))
		return
	}

	id, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Invalid user session"))
		return
	}

	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	user, err := uc.userService.UpdateProfile(id, req)
	if err != nil {
		c.JSON(userStatusCode(err), response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("User profile updated successfully", user))
}

func userStatusCode(err error) int {
	switch {
	case errors.Is(err, service.ErrInvalidUserID), errors.Is(err, service.ErrInvalidEmail):
		return http.StatusBadRequest
	case errors.Is(err, service.ErrUserNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
