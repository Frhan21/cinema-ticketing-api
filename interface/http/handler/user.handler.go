package handler

import (
	userdto "cinema-ticketing-api/app/user/dto"
	userservice "cinema-ticketing-api/app/user/service"
	"cinema-ticketing-api/entities"
	"cinema-ticketing-api/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService userservice.UserService
}

func NewUserController(userService userservice.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func toUserResponse(u *entities.User) *userdto.UserResponse {
	return &userdto.UserResponse{
		ID:    u.ID.String(),
		Name:  u.Name,
		Email: u.Email,
		Role:  string(u.Role),
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
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("User profile fetched successfully", toUserResponse(user)))
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

	var req userdto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	user, err := uc.userService.UpdateProfile(id, req)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("User profile updated successfully", toUserResponse(user)))
}
