package controller

import (
	"cinema-ticketing-api/models"
	"cinema-ticketing-api/request"
	"cinema-ticketing-api/response"
	"cinema-ticketing-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type studioController struct {
	studioService service.StudioService
}

func NewStudioController(studioService service.StudioService) *studioController {
	return &studioController{studioService: studioService}
}

func (sc *studioController) GetAll(c *gin.Context) {
	studios, err := sc.studioService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Success get studios", studios))
}

func (sc *studioController) GetByID(c *gin.Context) {
	id := c.Param("id")
	studio, err := sc.studioService.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Success get studio", studio))
}

func (sc *studioController) Create(c *gin.Context) {
	var req request.StudioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	studio := models.Studio{
		Name:       req.Name,
		Capacity:   req.Capacity,
		Facilities: req.Facilities,
	}

	if err := sc.studioService.Create(&studio); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Success create studio", studio))
}

func (sc *studioController) Update(c *gin.Context) {
	id := c.Param("id")
	
	// Cek apakah studio exists terlebih dahulu
	existingStudio, err := sc.studioService.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse("Studio not found"))
		return
	}

	var req request.UpdateStudioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	if req.Name != nil {
		existingStudio.Name = *req.Name
	}
	if req.Capacity != nil {
		existingStudio.Capacity = *req.Capacity
	}
	if req.Facilities != nil {
		existingStudio.Facilities = *req.Facilities
	}

	if err := sc.studioService.Update(existingStudio); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Success update studio", existingStudio))
}

func (sc *studioController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := sc.studioService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Success delete studio", nil))
}
