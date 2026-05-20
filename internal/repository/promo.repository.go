package repository

import (
	"cinema-ticketing-api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PromoRepository interface {
	Create(promo *model.Promo) error
	FindAll() ([]model.Promo, error)
	FindByID(id uuid.UUID) (*model.Promo, error)
	FindByCode(code string) (*model.Promo, error)
	Update(promo *model.Promo) error
	IncrementUsage(id uuid.UUID) error
	Delete(id uuid.UUID) error
}

type promoRepository struct {
	db *gorm.DB
}

func NewPromoRepository(db *gorm.DB) PromoRepository {
	return &promoRepository{db: db}
}

func (r *promoRepository) Create(promo *model.Promo) error {
	return r.db.Create(promo).Error
}

func (r *promoRepository) FindAll() ([]model.Promo, error) {
	var promos []model.Promo
	err := r.db.Order("created_at DESC").Find(&promos).Error
	return promos, err
}

func (r *promoRepository) FindByID(id uuid.UUID) (*model.Promo, error) {
	var promo model.Promo
	err := r.db.Where("id = ?", id).First(&promo).Error
	return &promo, err
}

func (r *promoRepository) FindByCode(code string) (*model.Promo, error) {
	var promo model.Promo
	err := r.db.Where("code = ?", code).First(&promo).Error
	return &promo, err
}

func (r *promoRepository) Update(promo *model.Promo) error {
	return r.db.Save(promo).Error
}

// IncrementUsage menambah used_count sebesar 1 secara atomic.
func (r *promoRepository) IncrementUsage(id uuid.UUID) error {
	return r.db.Model(&model.Promo{}).Where("id = ?", id).
		UpdateColumn("used_count", gorm.Expr("used_count + 1")).Error
}

func (r *promoRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&model.Promo{}).Error
}
