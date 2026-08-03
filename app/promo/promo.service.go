package promo

import (
	"cinema-ticketing-api/entities"
	"cinema-ticketing-api/pkg/apperror"
	"time"

	"github.com/google/uuid"
)

type PromoService interface {
	Create(req CreatePromoRequest) (*entities.Promo, error)
	GetAll() ([]entities.Promo, error)
	GetByID(id uuid.UUID) (*entities.Promo, error)
	ValidateCode(code string, totalPrice float64) (*entities.Promo, float64, error)
	Update(id uuid.UUID, req UpdatePromoRequest) (*entities.Promo, error)
	Delete(id uuid.UUID) error
}

type promoService struct {
	promoRepo PromoRepository
}

func NewPromoService(promoRepo PromoRepository) PromoService {
	return &promoService{promoRepo: promoRepo}
}

func (s *promoService) Create(req CreatePromoRequest) (*entities.Promo, error) {
	// Cek duplikat kode
	if existing, _ := s.promoRepo.FindByCode(req.Code); existing != nil {
		return nil, apperror.NewBadRequestError("promo code already exists")
	}

	if req.EndDate.Before(req.StartDate) {
		return nil, apperror.NewBadRequestError("end_date must be after start_date")
	}

	promo := entities.Promo{
		ID:          uuid.New(),
		Code:        req.Code,
		Description: req.Description,
		Discount:    req.Discount,
		MaxUsage:    req.MaxUsage,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		IsActive:    true,
	}

	if err := s.promoRepo.Create(&promo); err != nil {
		return nil, apperror.NewInternalServerError("failed to create promo")
	}

	return &promo, nil
}

func (s *promoService) GetAll() ([]entities.Promo, error) {
	promos, err := s.promoRepo.FindAll()
	if err != nil {
		return nil, err
	}

	return promos, nil
}

func (s *promoService) GetByID(id uuid.UUID) (*entities.Promo, error) {
	promo, err := s.promoRepo.FindByID(id)
	if err != nil {
		return nil, apperror.NewNotFoundError("promo not found")
	}
	return promo, nil
}

// ValidateCode memvalidasi kode promo dan menghitung harga setelah diskon.
func (s *promoService) ValidateCode(code string, totalPrice float64) (*entities.Promo, float64, error) {
	promo, err := s.promoRepo.FindByCode(code)
	if err != nil {
		return nil, 0, apperror.NewNotFoundError("promo code not found")
	}

	now := time.Now()
	if !promo.IsActive {
		return nil, 0, apperror.NewBadRequestError("promo is not active")
	}
	if now.Before(promo.StartDate) {
		return nil, 0, apperror.NewBadRequestError("promo has not started yet")
	}
	if now.After(promo.EndDate) {
		return nil, 0, apperror.NewBadRequestError("promo has expired")
	}
	if promo.MaxUsage > 0 && promo.UsedCount >= promo.MaxUsage {
		return nil, 0, apperror.NewBadRequestError("promo usage limit reached")
	}

	discountAmount := totalPrice * (promo.Discount / 100)
	discountedPrice := totalPrice - discountAmount

	return promo, discountedPrice, nil
}

func (s *promoService) Update(id uuid.UUID, req UpdatePromoRequest) (*entities.Promo, error) {
	promo, err := s.promoRepo.FindByID(id)
	if err != nil {
		return nil, apperror.NewNotFoundError("promo not found")
	}

	if req.Description != "" {
		promo.Description = req.Description
	}
	if req.Discount > 0 {
		promo.Discount = req.Discount
	}
	if req.MaxUsage >= 0 {
		promo.MaxUsage = req.MaxUsage
	}
	if !req.StartDate.IsZero() {
		promo.StartDate = req.StartDate
	}
	if !req.EndDate.IsZero() {
		promo.EndDate = req.EndDate
	}
	if req.IsActive != nil {
		promo.IsActive = *req.IsActive
	}

	if err := s.promoRepo.Update(promo); err != nil {
		return nil, apperror.NewInternalServerError("failed to update promo")
	}

	return promo, nil
}

func (s *promoService) Delete(id uuid.UUID) error {
	if _, err := s.promoRepo.FindByID(id); err != nil {
		return apperror.NewNotFoundError("promo not found")
	}
	return s.promoRepo.Delete(id)
}
