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

// Create godoc
// @Summary      Buat promo baru
// @Description  Admin dapat membuat kode promo diskon dengan persentase tertentu dan periode berlaku.
// @Tags         Promo
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  request.CreatePromoRequest  true  "Data promo"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Router       /promo [post]
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
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Promos retrieved successfully", promos))
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
		c.JSON(http.StatusNotFound, response.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Promo retrieved successfully", promo))
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
	result, err := p.promoService.ValidateCode(code, req.TotalPrice)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Promo is valid", result))
}

// Update godoc
// @Summary      Update promo
// @Description  Admin dapat memperbarui data promo yang sudah ada.
// @Tags         Promo
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                    true  "Promo ID (UUID)"
// @Param        body  body  request.UpdatePromoRequest  true  "Data promo yang diperbarui"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /promo/{id} [put]
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
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Promo deleted successfully", nil))
}
