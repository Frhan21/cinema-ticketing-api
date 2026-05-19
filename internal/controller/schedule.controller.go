package controller

import (
	"cinema-ticketing-api/internal/request"
	"cinema-ticketing-api/internal/response"
	"cinema-ticketing-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ScheduleController interface {
	FindAll(c *gin.Context)
	FindByID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type scheduleController struct {
	scheduleService service.ScheduleService
}

// Create implements [ScheduleController].
func (s *scheduleController) Create(c *gin.Context) {
	var req request.CreateSchedule

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(err.Error())
		return
	}

	schedule, err := s.scheduleService.Create(req)
	if err != nil {
		response.ErrorResponse(err.Error())
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("Schedule created successfully", &response.CreateScheduleResponse{
		ID:        schedule.ID,
		MovieID:   schedule.MovieID,
		StudioID:  schedule.StudioID,
		StartTime: schedule.StartTime.Format("2006-01-02T15:04:05Z"),
		EndTime:   schedule.EndTime.Format("2006-01-02T15:04:05Z"),
		Price:     schedule.Price,
	}))
}

// Delete implements [ScheduleController].
func (s *scheduleController) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ErrorResponse("Schedule ID is required")
		return
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}
	err = s.scheduleService.Delete(parsedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Schedule deleted successfully", nil))
}

// FindAll implements [ScheduleController].
func (s *scheduleController) FindAll(c *gin.Context) {
	var req request.ScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}
	schedules, err := s.scheduleService.FindAll(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Schedules fetched successfully", schedules))
}

// FindByID implements [ScheduleController].
func (s *scheduleController) FindByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ErrorResponse("Schedule ID is required")
		return
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}
	schedule, err := s.scheduleService.FindByID(parsedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Schedule fetched successfully", schedule))
}

// Update implements [ScheduleController].
func (s *scheduleController) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ErrorResponse("Schedule ID is required")
		return
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}
	var req request.UpdateSchedule
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}
	schedule, err := s.scheduleService.Update(parsedID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Schedule updated successfully", &response.UpdateScheduleResponse{
		ID:        schedule.ID,
		MovieID:   schedule.MovieID,
		StudioID:  schedule.StudioID,
		StartTime: schedule.StartTime.Format("2006-01-02T15:04:05Z"),
		EndTime:   schedule.EndTime.Format("2006-01-02T15:04:05Z"),
		Price:     schedule.Price,
	}))
}

func NewScheduleController(scheduleService service.ScheduleService) ScheduleController {
	return &scheduleController{scheduleService: scheduleService}
}
