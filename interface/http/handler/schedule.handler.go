package handler

import (
	scheduledto "cinema-ticketing-api/app/schedule/dto"
	scheduleservice "cinema-ticketing-api/app/schedule/service"
	"cinema-ticketing-api/pkg/apperror"
	"cinema-ticketing-api/pkg/pagination"
	"cinema-ticketing-api/pkg/utils"
	"cinema-ticketing-api/request"
	"cinema-ticketing-api/response"
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
	scheduleService scheduleservice.ScheduleService
}

// Create implements [ScheduleController].
func (s *scheduleController) Create(c *gin.Context) {
	var req scheduledto.CreateSchedule

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	schedule, err := s.scheduleService.Create(req)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("Schedule created successfully", &scheduledto.CreateScheduleResponse{
		ID:        schedule.ID,
		MovieID:   schedule.MovieID,
		StudioID:  schedule.StudioID,
		StartTime: utils.FormatTime(schedule.StartTime),
		EndTime:   utils.FormatTime(schedule.EndTime),
		Price:     schedule.Price,
	}))
}

// Delete implements [ScheduleController].
func (s *scheduleController) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.Error(apperror.NewBadRequestError("Schedule ID is required"))
		c.Abort()
		return
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	err = s.scheduleService.Delete(parsedID)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Schedule deleted successfully", nil))
}

// FindAll implements [ScheduleController].
func (s *scheduleController) FindAll(c *gin.Context) {
	var req request.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	schedules, count, err := s.scheduleService.FindAll(req)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	meta := &response.Meta{
		Page:      req.Page,
		PerPage:   req.PerPage,
		TotalData: int(count),
		TotalPage: pagination.GetTotalPage(count, req.PerPage),
	}

	var scheduleResponses []scheduledto.ScheduleResponse
	for _, sch := range schedules {
		scheduleResponses = append(scheduleResponses, scheduledto.ScheduleResponse{
			ID:        sch.ID,
			MovieID:   sch.MovieID,
			StudioID:  sch.StudioID,
			StartTime: utils.FormatTime(sch.StartTime),
			EndTime:   utils.FormatTime(sch.EndTime),
			Price:     sch.Price,
		})
	}

	c.JSON(http.StatusOK, response.SuccessResponseWithMeta("Schedules fetched successfully", scheduleResponses, meta))
}

// FindByID implements [ScheduleController].
func (s *scheduleController) FindByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.Error(apperror.NewBadRequestError("Schedule ID is required"))
		c.Abort()
		return
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	schedule, err := s.scheduleService.FindByID(parsedID)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Schedule fetched successfully", &scheduledto.ScheduleResponse{
		ID:        schedule.ID,
		MovieID:   schedule.MovieID,
		StudioID:  schedule.StudioID,
		StartTime: utils.FormatTime(schedule.StartTime),
		EndTime:   utils.FormatTime(schedule.EndTime),
		Price:     schedule.Price,
	}))
}

// Update implements [ScheduleController].
func (s *scheduleController) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.Error(apperror.NewBadRequestError("Schedule ID is required"))
		c.Abort()
		return
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	var req scheduledto.UpdateSchedule
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	schedule, err := s.scheduleService.Update(parsedID, req)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Schedule updated successfully", &scheduledto.UpdateScheduleResponse{
		ID:        schedule.ID,
		MovieID:   schedule.MovieID,
		StudioID:  schedule.StudioID,
		StartTime: utils.FormatTime(schedule.StartTime),
		EndTime:   utils.FormatTime(schedule.EndTime),
		Price:     schedule.Price,
	}))
}

func NewScheduleController(scheduleService scheduleservice.ScheduleService) ScheduleController {
	return &scheduleController{scheduleService: scheduleService}
}
