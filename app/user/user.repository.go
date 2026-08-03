package user

import (
	"cinema-ticketing-api/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*entities.User, error)
	FindByID(id uuid.UUID) (*entities.User, error)
	Update(user *entities.User) error
	Create(user *entities.User) error
	Delete(id uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (us *userRepository) Create(user *entities.User) error {
	return us.db.Table("users").Create(user).Error
}

func (us *userRepository) FindByEmail(email string) (*entities.User, error) {

	user := &entities.User{}
	err := us.db.Table("users").Where("email = ?", email).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *userRepository) FindByID(id uuid.UUID) (*entities.User, error) {
	user := &entities.User{}

	err := us.db.Table("users").Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *userRepository) Update(user *entities.User) error {
	return us.db.Table("users").Save(user).Error
}

func (us *userRepository) Delete(id uuid.UUID) error {
	return us.db.Table("users").Delete(&entities.User{}, "id = ?", id).Error
}
