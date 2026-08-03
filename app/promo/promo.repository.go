package promo

import (
	"cinema-ticketing-api/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PromoRepository interface {
	Create(promo *entities.Promo) error
	FindAll() ([]entities.Promo, error)
	FindByID(id uuid.UUID) (*entities.Promo, error)
	FindByCode(code string) (*entities.Promo, error)
	Update(promo *entities.Promo) error
	IncrementUsage(id uuid.UUID) error
	Delete(id uuid.UUID) error
}

type promoRepository struct {
	db *gorm.DB
}

func NewPromoRepository(db *gorm.DB) PromoRepository {
	return &promoRepository{db: db}
}

func (r *promoRepository) Create(promo *entities.Promo) error {
	return r.db.Create(promo).Error
}

func (r *promoRepository) FindAll() ([]entities.Promo, error) {
	var promos []entities.Promo
	err := r.db.Order("created_at DESC").Find(&promos).Error
	return promos, err
}

func (r *promoRepository) FindByID(id uuid.UUID) (*entities.Promo, error) {
	var promo entities.Promo
	err := r.db.Where("id = ?", id).First(&promo).Error
	return &promo, err
}

func (r *promoRepository) FindByCode(code string) (*entities.Promo, error) {
	var promo entities.Promo
	err := r.db.Where("code = ?", code).First(&promo).Error
	return &promo, err
}

func (r *promoRepository) Update(promo *entities.Promo) error {
	return r.db.Save(promo).Error
}

// IncrementUsage menambah used_count sebesar 1 secara atomic.
func (r *promoRepository) IncrementUsage(id uuid.UUID) error {
	return r.db.Model(&entities.Promo{}).Where("id = ?", id).
		UpdateColumn("used_count", gorm.Expr("used_count + 1")).Error
}

func (r *promoRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&entities.Promo{}).Error
}
