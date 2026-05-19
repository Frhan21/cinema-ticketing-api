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
	"github.com/google/uuid"
)

func seatStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch {
	case errors.Is(err, service.ErrSeatNotFound):
		return http.StatusNotFound
	case errors.Is(err, service.ErrInvalidSeatID):
		return http.StatusBadRequest
	case errors.Is(err, service.ErrInvalidStudioID):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

type seatController struct {
	seatService service.SeatService
}

func (s *seatController) Create(c *gin.Context) {
	var req request.SeatRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	studioID, err := uuid.Parse(req.StudioID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid Studio ID format"))
		return
	}

	isAvailable := true
	if req.IsAvailable != nil {
		isAvailable = *req.IsAvailable
	}

	seat := &model.Seat{
		StudioID:    studioID,
		SeatNumber:  req.SeatNumber,
		IsAvailable: isAvailable,
	}

	err = s.seatService.Create(seat)
	if err != nil {
		c.JSON(seatStatusCode(err), response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("Seat created successfully", seat))
}

func (s *seatController) FindAll(c *gin.Context) {
	var req request.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	seats, count, err := s.seatService.FindAll(req)
	if err != nil {
		c.JSON(seatStatusCode(err), response.ErrorResponse(err.Error()))
		return
	}

	meta := &response.Meta{
		Page:      req.Page,
		PerPage:   req.PerPage,
		TotalData: int(count),
		TotalPage: pagination.GetTotalPage(count, req.PerPage),
	}

	c.JSON(http.StatusOK, response.SuccessResponseWithMeta("Seats fetched successfully", seats, meta))
}

func (s *seatController) FindByID(c *gin.Context) {
	id := c.Param("id")

	seat, err := s.seatService.FindByID(id)
	if err != nil {
		c.JSON(seatStatusCode(err), response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Seat fetched successfully", seat))
}

func (s *seatController) FindByStudioID(c *gin.Context) {
	studioID := c.Param("studio_id")

	seats, err := s.seatService.FindByStudioID(studioID)
	if err != nil {
		c.JSON(seatStatusCode(err), response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Seats fetched successfully", seats))
}

func (s *seatController) Update(c *gin.Context) {
	id := c.Param("id")

	seat, err := s.seatService.FindByID(id)
	if err != nil {
		c.JSON(seatStatusCode(err), response.ErrorResponse(err.Error()))
		return
	}

	var req request.UpdateSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	if req.SeatNumber != "" {
		seat.SeatNumber = req.SeatNumber
	}

	if req.StudioID != "" {
		studioID, err := uuid.Parse(req.StudioID)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid Studio ID format"))
			return
		}
		seat.StudioID = studioID
	}

	if req.IsAvailable != nil {
		seat.IsAvailable = *req.IsAvailable
	}

	err = s.seatService.Update(seat)
	if err != nil {
		c.JSON(seatStatusCode(err), response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Seat updated successfully", seat))
}

func (s *seatController) Delete(c *gin.Context) {
	id := c.Param("id")

	err := s.seatService.Delete(id)
	if err != nil {
		c.JSON(seatStatusCode(err), response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Seat deleted successfully", nil))
}

type SeatController interface {
	FindAll(c *gin.Context)
	FindByID(c *gin.Context)
	FindByStudioID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

func NewSeatController(seatService service.SeatService) SeatController {
	return &seatController{seatService: seatService}
}
