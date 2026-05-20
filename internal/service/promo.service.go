package service

import (
	"cinema-ticketing-api/internal/model"
	"cinema-ticketing-api/internal/repository"
	"cinema-ticketing-api/internal/request"
	"cinema-ticketing-api/internal/response"
	"errors"
	"time"

	"github.com/google/uuid"
)

type PromoService interface {
	Create(req request.CreatePromoRequest) (*response.PromoResponse, error)
	GetAll() ([]response.PromoResponse, error)
	GetByID(id uuid.UUID) (*response.PromoResponse, error)
	ValidateCode(code string, totalPrice float64) (*response.ValidatePromoResponse, error)
	Update(id uuid.UUID, req request.UpdatePromoRequest) (*response.PromoResponse, error)
	Delete(id uuid.UUID) error
}

type promoService struct {
	promoRepo repository.PromoRepository
}

func NewPromoService(promoRepo repository.PromoRepository) PromoService {
	return &promoService{promoRepo: promoRepo}
}

func (s *promoService) Create(req request.CreatePromoRequest) (*response.PromoResponse, error) {
	// Cek duplikat kode
	if existing, _ := s.promoRepo.FindByCode(req.Code); existing != nil {
		return nil, errors.New("promo code already exists")
	}

	if req.EndDate.Before(req.StartDate) {
		return nil, errors.New("end_date must be after start_date")
	}

	promo := model.Promo{
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
		return nil, errors.New("failed to create promo")
	}

	return toPromoResponse(promo), nil
}

func (s *promoService) GetAll() ([]response.PromoResponse, error) {
	promos, err := s.promoRepo.FindAll()
	if err != nil {
		return nil, err
	}
	var result []response.PromoResponse
	for _, p := range promos {
		result = append(result, *toPromoResponse(p))
	}
	return result, nil
}

func (s *promoService) GetByID(id uuid.UUID) (*response.PromoResponse, error) {
	promo, err := s.promoRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("promo not found")
	}
	return toPromoResponse(*promo), nil
}

// ValidateCode memvalidasi kode promo dan menghitung harga setelah diskon.
func (s *promoService) ValidateCode(code string, totalPrice float64) (*response.ValidatePromoResponse, error) {
	promo, err := s.promoRepo.FindByCode(code)
	if err != nil {
		return nil, errors.New("promo code not found")
	}

	now := time.Now()
	if !promo.IsActive {
		return nil, errors.New("promo is not active")
	}
	if now.Before(promo.StartDate) {
		return nil, errors.New("promo has not started yet")
	}
	if now.After(promo.EndDate) {
		return nil, errors.New("promo has expired")
	}
	if promo.MaxUsage > 0 && promo.UsedCount >= promo.MaxUsage {
		return nil, errors.New("promo usage limit reached")
	}

	discountAmount := totalPrice * (promo.Discount / 100)
	discountedPrice := totalPrice - discountAmount

	return &response.ValidatePromoResponse{
		Code:            promo.Code,
		Discount:        promo.Discount,
		OriginalPrice:   totalPrice,
		DiscountedPrice: discountedPrice,
	}, nil
}

func (s *promoService) Update(id uuid.UUID, req request.UpdatePromoRequest) (*response.PromoResponse, error) {
	promo, err := s.promoRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("promo not found")
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
		return nil, errors.New("failed to update promo")
	}

	return toPromoResponse(*promo), nil
}

func (s *promoService) Delete(id uuid.UUID) error {
	if _, err := s.promoRepo.FindByID(id); err != nil {
		return errors.New("promo not found")
	}
	return s.promoRepo.Delete(id)
}

// toPromoResponse mengkonversi model.Promo ke response DTO.
func toPromoResponse(p model.Promo) *response.PromoResponse {
	return &response.PromoResponse{
		ID:          p.ID,
		Code:        p.Code,
		Description: p.Description,
		Discount:    p.Discount,
		MaxUsage:    p.MaxUsage,
		UsedCount:   p.UsedCount,
		StartDate:   p.StartDate,
		EndDate:     p.EndDate,
		IsActive:    p.IsActive,
	}
}
