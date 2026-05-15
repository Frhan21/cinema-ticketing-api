package service

import (
	"cinema-ticketing-api/models"
	"cinema-ticketing-api/repository"
	"cinema-ticketing-api/request"
	"cinema-ticketing-api/response"
	"cinema-ticketing-api/utils"
	"errors"

	"gorm.io/gorm"
)

type AuthService interface {
	Register(req request.RegisterRequest) (*response.AuthResponse, error)
	Login(req request.LoginRequest) (*response.AuthResponse, error)
}

type authService struct {
	userRepository repository.UserRepository
}

// Login implements [AuthService].
func (a *authService) Login(req request.LoginRequest) (*response.AuthResponse, error) {
	existingUser, err := a.userRepository.FindByEmail(req.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Invalid email or password")
		}
		return nil, err
	}

	if !utils.CheckPasswordHash(req.Password, existingUser.Password) {
		return nil, errors.New("Invalid email or password")
	}

	token, err := utils.GenerateJwt(existingUser.ID, string(existingUser.Role))

	if err != nil {
		return nil, err
	}

	user := &response.AuthResponse{
		ID:    existingUser.ID.String(),
		Name:  existingUser.Name,
		Email: existingUser.Email,
		Role:  string(existingUser.Role),
		Token: token,
	}

	return user, nil
}

// Register implements [AuthService].
func (a *authService) Register(req request.RegisterRequest) (*response.AuthResponse, error) {

	existUser, err := a.userRepository.FindByEmail(req.Email)
	if err == nil && existUser != nil {
		return nil, errors.New("Email already exists")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if req.Password != req.ConfirmPassword {
		return nil, errors.New("Password and confirm password do not match")
	}

	hashsedPassword, err := utils.HashPassword(req.Password)

	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashsedPassword,
		Role:     models.UserRole,
	}

	err = a.userRepository.Create(&user)

	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateJwt(user.ID, string(user.Role))

	if err != nil {
		return nil, err
	}

	res := &response.AuthResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
		Token: token,
	}

	return res, nil
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepository: userRepo,
	}
}
