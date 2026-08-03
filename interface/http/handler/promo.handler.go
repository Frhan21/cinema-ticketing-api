package handler

import (
	promopkg "cinema-ticketing-api/app/promo"
	"cinema-ticketing-api/entities"
	"cinema-ticketing-api/response"
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
	promoService promopkg.PromoService
}

func NewPromoController(promoService promopkg.PromoService) PromoController {
	return &promoController{promoService: promoService}
}

// toPromoResponse mengkonversi model Promo ke response DTO.
func toPromoResponse(p entities.Promo) *promopkg.PromoResponse {
	return &promopkg.PromoResponse{
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

// Create godoc
// @Summary      Buat promo baru
// @Description  Admin dapat membuat kode promo diskon dengan persentase tertentu dan periode berlaku.
// @Tags         Promo
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  promopkg.CreatePromoRequest  true  "Data promo"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Router       /promo [post]
func (p *promoController) Create(c *gin.Context) {
	var req promopkg.CreatePromoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	promo, err := p.promoService.Create(req)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusCreated, response.SuccessResponse("Promo created successfully", toPromoResponse(*promo)))
}

// GetAll godoc
// @Summary      List semua promo
// @Description  Admin dapat melihat semua kode promo yang tersedia.
// @Tags         Promo
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /promo [get]
func (p *promoController) GetAll(c *gin.Context) {
	promos, err := p.promoService.GetAll()
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	var res []promopkg.PromoResponse
	for _, promo := range promos {
		res = append(res, *toPromoResponse(promo))
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Promos retrieved successfully", res))
}

// GetByID godoc
// @Summary      Detail promo
// @Description  Mengambil detail satu kode promo berdasarkan ID-nya.
// @Tags         Promo
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Promo ID (UUID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /promo/{id} [get]
func (p *promoController) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid promo ID"))
		return
	}
	promo, err := p.promoService.GetByID(id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Promo retrieved successfully", toPromoResponse(*promo)))
}

// ValidateCode godoc
// @Summary      Validasi kode promo
// @Description  User dapat mengecek apakah kode promo valid dan melihat preview harga setelah diskon sebelum melakukan booking.
// @Tags         Promo
// @Produce      json
// @Security     BearerAuth
// @Param        code         path   string   true   "Kode promo"
// @Param        total_price  query  number   true   "Harga total sebelum diskon"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /promo/validate/{code} [get]
func (p *promoController) ValidateCode(c *gin.Context) {
	code := c.Param("code")
	var req struct {
		TotalPrice float64 `form:"total_price" binding:"required,min=0"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("total_price query param is required"))
		return
	}
	promo, discountedPrice, err := p.promoService.ValidateCode(code, req.TotalPrice)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	res := promopkg.ValidatePromoResponse{
		Code:            promo.Code,
		Discount:        promo.Discount,
		OriginalPrice:   req.TotalPrice,
		DiscountedPrice: discountedPrice,
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Promo is valid", res))
}

// Update godoc
// @Summary      Update promo
// @Description  Admin dapat memperbarui data promo yang sudah ada.
// @Tags         Promo
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                    true  "Promo ID (UUID)"
// @Param        body  body  promopkg.UpdatePromoRequest         true  "Data promo yang diperbarui"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /promo/{id} [put]
func (p *promoController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid promo ID"))
		return
	}
	var req promopkg.UpdatePromoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}
	promo, err := p.promoService.Update(id, req)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Promo updated successfully", toPromoResponse(*promo)))
}

// Delete godoc
// @Summary      Hapus promo
// @Description  Admin dapat menghapus kode promo yang sudah tidak digunakan.
// @Tags         Promo
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Promo ID (UUID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /promo/{id} [delete]
func (p *promoController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid promo ID"))
		return
	}
	if err := p.promoService.Delete(id); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Promo deleted successfully", nil))
}
