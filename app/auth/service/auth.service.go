package service

import (
	authdto "cinema-ticketing-api/app/auth/dto"
	userrepository "cinema-ticketing-api/app/user/repository"
	"cinema-ticketing-api/config"
	"cinema-ticketing-api/entities"
	"cinema-ticketing-api/pkg/apperror"
	"cinema-ticketing-api/pkg/jwt"
	"cinema-ticketing-api/pkg/password"
	"errors"

	"gorm.io/gorm"
)

type AuthService interface {
	Register(req authdto.RegisterRequest) (*entities.User, string, error)
	Login(req authdto.LoginRequest) (*entities.User, string, error)
}

type authService struct {
	userRepository userrepository.UserRepository
	jwtConfig      config.JWTConfig
}

// Login implements [AuthService].
func (a *authService) Login(req authdto.LoginRequest) (*entities.User, string, error) {
	existingUser, err := a.userRepository.FindByEmail(req.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", apperror.NewUnauthorizedError("invalid email or password")
		}
		return nil, "", err
	}

	if !password.CheckPasswordHash(req.Password, existingUser.Password) {
		return nil, "", apperror.NewUnauthorizedError("invalid email or password")
	}

	token, err := jwt.GenerateJwt(existingUser.ID, string(existingUser.Role), a.jwtConfig)
	if err != nil {
		return nil, "", err
	}

	return existingUser, token, nil
}

// Register implements [AuthService].
func (a *authService) Register(req authdto.RegisterRequest) (*entities.User, string, error) {
	existUser, err := a.userRepository.FindByEmail(req.Email)
	if err == nil && existUser != nil {
		return nil, "", apperror.NewBadRequestError("email already exists")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", err
	}

	if req.Password != req.ConfirmPassword {
		return nil, "", apperror.NewBadRequestError("password and confirm password do not match")
	}

	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		return nil, "", err
	}

	user := entities.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     entities.RoleUser,
	}

	if err = a.userRepository.Create(&user); err != nil {
		return nil, "", err
	}

	token, err := jwt.GenerateJwt(user.ID, string(user.Role), a.jwtConfig)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func NewAuthService(userRepo userrepository.UserRepository, jwtCfg config.JWTConfig) AuthService {
	return &authService{
		userRepository: userRepo,
		jwtConfig:      jwtCfg,
	}
}
