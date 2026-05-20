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
	return &AuthController{authService: authService}
}

// Register godoc
// @Summary      Register user baru
// @Description  Membuat akun pengguna baru dengan role 'user'.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body request.RegisterRequest true "Data registrasi"
// @Success      201  {object}  map[string]interface{}  "User berhasil didaftarkan"
// @Failure      400  {object}  map[string]interface{}  "Request tidak valid"
// @Router       /auth/register [post]
func (a *AuthController) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	result, err := a.authService.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("User registered successfully", result))
}

// Login godoc
// @Summary      Login user
// @Description  Autentikasi user dan mendapatkan JWT token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body request.LoginRequest true "Kredensial login"
// @Success      200  {object}  map[string]interface{}  "Login berhasil, token tersedia"
// @Failure      400  {object}  map[string]interface{}  "Request tidak valid"
// @Failure      401  {object}  map[string]interface{}  "Email atau password salah"
// @Router       /auth/login [post]
func (a *AuthController) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	result, err := a.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Login successful", result))
}
