package handler

import (
	"cinema-ticketing-api/app/auth"
	"cinema-ticketing-api/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService auth.AuthService
}

func NewAuthController(authService auth.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Register godoc
// @Summary      Register user baru
// @Description  Membuat akun pengguna baru dengan role 'user'.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body auth.RegisterRequest true "Data registrasi"
// @Success      201  {object}  map[string]interface{}  "User berhasil didaftarkan"
// @Failure      400  {object}  map[string]interface{}  "Request tidak valid"
// @Router       /auth/register [post]
func (a *AuthController) Register(c *gin.Context) {
	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	user, token, err := a.authService.Register(req)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	res := auth.AuthResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
		Token: token,
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("User registered successfully", res))
}

// Login godoc
// @Summary      Login user
// @Description  Autentikasi user dan mendapatkan JWT token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body auth.LoginRequest true "Kredensial login"
// @Success      200  {object}  map[string]interface{}  "Login berhasil, token tersedia"
// @Failure      400  {object}  map[string]interface{}  "Request tidak valid"
// @Failure      401  {object}  map[string]interface{}  "Email atau password salah"
// @Router       /auth/login [post]
func (a *AuthController) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	user, token, err := a.authService.Login(req)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	res := auth.AuthResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
		Token: token,
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Login successful", res))
}
