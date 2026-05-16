package service

import (
	"cinema-ticketing-api/internal/model"
	"cinema-ticketing-api/internal/repository"
	"cinema-ticketing-api/internal/request"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrInvalidUserID = errors.New("invalid user id")
	ErrInvalidEmail  = errors.New("email already used by another user")
)

type UserService interface {
	GetProfile(id string) (*model.User, error)
	UpdateProfile(id string, input request.UpdateUserRequest) (*model.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) GetProfile(id string) (*model.User, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidUserID
	}

	user, err := s.userRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateProfile(id string, input request.UpdateUserRequest) (*model.User, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidUserID
	}

	user, err := s.userRepository.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if input.Name != nil {
		user.Name = strings.TrimSpace(*input.Name)
	}

	if input.Email != nil {
		email := strings.TrimSpace(*input.Email)
		existingUser, err := s.userRepository.FindByEmail(email)
		if err == nil && existingUser != nil && existingUser.ID != user.ID {
			return nil, ErrInvalidEmail
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		user.Email = email
	}

	if err = s.userRepository.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
