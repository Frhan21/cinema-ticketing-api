package repository

import (
	"cinema-ticketing-api/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*models.User, error)
	FindByID(id string) (*models.User, error)
	Update(user *models.User) error
	Create(user *models.User) error
	Delete(id string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (us *userRepository) Create(user *models.User) error {
	return us.db.Table("users").Create(user).Error
}

func (us *userRepository) FindByEmail(email string) (*models.User, error) {

	user := &models.User{}
	err := us.db.Table("users").Where("email = ?", email).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *userRepository) FindByID(id string) (*models.User, error) {
	user := &models.User{}

	err := us.db.Table("users").Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *userRepository) Update(user *models.User) error {
	return us.db.Table("users").Save(user).Error
}

func (us *userRepository) Delete(id string) error {
	return us.db.Table("users").Delete(&models.User{}, "id = ?", id).Error
}
