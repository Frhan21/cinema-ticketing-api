package repository

import (
	"cinema-ticketing-api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*model.User, error)
	FindByID(id uuid.UUID) (*model.User, error)
	Update(user *model.User) error
	Create(user *model.User) error
	Delete(id uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (us *userRepository) Create(user *model.User) error {
	return us.db.Table("users").Create(user).Error
}

func (us *userRepository) FindByEmail(email string) (*model.User, error) {

	user := &model.User{}
	err := us.db.Table("users").Where("email = ?", email).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *userRepository) FindByID(id uuid.UUID) (*model.User, error) {
	user := &model.User{}

	err := us.db.Table("users").Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *userRepository) Update(user *model.User) error {
	return us.db.Table("users").Save(user).Error
}

func (us *userRepository) Delete(id uuid.UUID) error {
	return us.db.Table("users").Delete(&model.User{}, "id = ?", id).Error
}
