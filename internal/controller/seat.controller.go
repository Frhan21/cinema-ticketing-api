package controller

import (
	"cinema-ticketing-api/internal/model"
	"cinema-ticketing-api/internal/request"
	"cinema-ticketing-api/internal/response"
	"cinema-ticketing-api/internal/service"
	"cinema-ticketing-api/pkg/pagination"
	"net/http"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type seatController struct {
	seatService service.SeatService
}

func toSeatResponse(s *model.Seat) *response.SeatResponse {
	return &response.SeatResponse{
		ID:          s.ID.String(),
		StudioID:    s.StudioID.String(),
		SeatNumber:  s.SeatNumber,
		IsAvailable: s.IsAvailable,
		CreatedAt:   s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   s.UpdatedAt.Format(time.RFC3339),
	}
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
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("Seat created successfully", toSeatResponse(seat)))
}

func (s *seatController) FindAll(c *gin.Context) {
	var req request.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	seats, count, err := s.seatService.FindAll(req)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	var res []response.SeatResponse
	for _, seat := range seats {
		res = append(res, *toSeatResponse(&seat))
	}

	meta := &response.Meta{
		Page:      req.Page,
		PerPage:   req.PerPage,
		TotalData: int(count),
		TotalPage: pagination.GetTotalPage(count, req.PerPage),
	}

	c.JSON(http.StatusOK, response.SuccessResponseWithMeta("Seats fetched successfully", res, meta))
}

func (s *seatController) FindByID(c *gin.Context) {
	id := c.Param("id")

	seat, err := s.seatService.FindByID(id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Seat fetched successfully", toSeatResponse(seat)))
}

func (s *seatController) FindByStudioID(c *gin.Context) {
	studioID := c.Param("studio_id")

	seats, err := s.seatService.FindByStudioID(studioID)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	var res []response.SeatResponse
	for _, seat := range seats {
		res = append(res, *toSeatResponse(&seat))
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Seats fetched successfully", res))
}

func (s *seatController) Update(c *gin.Context) {
	id := c.Param("id")

	seat, err := s.seatService.FindByID(id)
	if err != nil {
		c.Error(err)
		c.Abort()
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
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Seat updated successfully", toSeatResponse(seat)))
}

func (s *seatController) Delete(c *gin.Context) {
	id := c.Param("id")

	err := s.seatService.Delete(id)
	if err != nil {
		c.Error(err)
		c.Abort()
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
