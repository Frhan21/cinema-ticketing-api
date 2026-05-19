package controller

import (
	"cinema-ticketing-api/internal/model"
	"cinema-ticketing-api/internal/request"
	"cinema-ticketing-api/internal/response"
	"cinema-ticketing-api/internal/service"
	"cinema-ticketing-api/pkg/pagination"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StudioController interface {
	GetAll(c *gin.Context)
	GetByID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type studioController struct {
	studioService service.StudioService
}

func NewStudioController(studioService service.StudioService) StudioController {
	return &studioController{studioService: studioService}
}

func (sc *studioController) GetAll(c *gin.Context) {
	var req request.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	studios, count, err := sc.studioService.FindAll(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	meta := &response.Meta{
		Page:      req.Page,
		PerPage:   req.PerPage,
		TotalData: int(count),
		TotalPage: pagination.GetTotalPage(count, req.PerPage),
	}

	c.JSON(http.StatusOK, response.SuccessResponseWithMeta("Success get studios", studios, meta))
}

func (sc *studioController) GetByID(c *gin.Context) {
	id := c.Param("id")
	studio, err := sc.studioService.FindByID(id)
	if err != nil {
		c.JSON(studioStatusCode(err), response.ErrorResponse(err.Error()))
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

	studio := model.Studio{
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
		c.JSON(studioStatusCode(err), response.ErrorResponse(err.Error()))
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
		c.JSON(studioStatusCode(err), response.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Success delete studio", nil))
}

func studioStatusCode(err error) int {
	switch {
	case errors.Is(err, service.ErrInvalidStudioID):
		return http.StatusBadRequest
	case errors.Is(err, service.ErrStudioNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
