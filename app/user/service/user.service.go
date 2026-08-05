package service

import (
	userdto "cinema-ticketing-api/app/user/dto"
	userrepository "cinema-ticketing-api/app/user/repository"
	"cinema-ticketing-api/entities"
	"cinema-ticketing-api/pkg/apperror"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	GetProfile(id string) (*entities.User, error)
	UpdateProfile(id string, input userdto.UpdateUserRequest) (*entities.User, error)
}

type userService struct {
	userRepository userrepository.UserRepository
}

func NewUserService(userRepository userrepository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) GetProfile(id string) (*entities.User, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid user id")
	}

	user, err := s.userRepository.FindByID(parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFoundError("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateProfile(id string, input userdto.UpdateUserRequest) (*entities.User, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid user id")
	}

	user, err := s.userRepository.FindByID(parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFoundError("user not found")
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
			return nil, apperror.NewBadRequestError("email already used by another user")
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
