package controller

import (
	"cinema-ticketing-api/internal/request"
	"cinema-ticketing-api/internal/response"
	"cinema-ticketing-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PromoController interface {
	Create(c *gin.Context)
	GetAll(c *gin.Context)
	GetByID(c *gin.Context)
	ValidateCode(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type promoController struct {
	promoService service.PromoService
}

func NewPromoController(promoService service.PromoService) PromoController {
	return &promoController{promoService: promoService}
}

// Create membuat promo baru (admin only).
func (p *promoController) Create(c *gin.Context) {
	var req request.CreatePromoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	promo, err := p.promoService.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("Promo created successfully", promo))
}

// GetAll mengembalikan semua promo (admin only).
func (p *promoController) GetAll(c *gin.Context) {
	promos, err := p.promoService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Promos retrieved successfully", promos))
}

// GetByID mengembalikan detail satu promo.
func (p *promoController) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid promo ID"))
		return
	}

	promo, err := p.promoService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Promo retrieved successfully", promo))
}

// ValidateCode memvalidasi kode promo dan menampilkan harga setelah diskon.
// User mengirim total harga via query param untuk preview diskon.
func (p *promoController) ValidateCode(c *gin.Context) {
	code := c.Param("code")

	// Ambil total_price dari query param untuk preview diskon
	var req struct {
		TotalPrice float64 `form:"total_price" binding:"required,min=0"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("total_price query param is required"))
		return
	}

	result, err := p.promoService.ValidateCode(code, req.TotalPrice)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Promo is valid", result))
}

// Update memperbarui data promo (admin only).
func (p *promoController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid promo ID"))
		return
	}

	var req request.UpdatePromoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	promo, err := p.promoService.Update(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Promo updated successfully", promo))
}

// Delete menghapus promo (admin only).
func (p *promoController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid promo ID"))
		return
	}

	if err := p.promoService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Promo deleted successfully", nil))
}
