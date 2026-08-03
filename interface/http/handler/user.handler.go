package handler

import (
	"cinema-ticketing-api/app/user"
	"cinema-ticketing-api/entities"
	"cinema-ticketing-api/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService user.UserService
}

func NewUserController(userService user.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func toUserResponse(u *entities.User) *user.UserResponse {
	return &user.UserResponse{
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

	var req user.UpdateUserRequest
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
