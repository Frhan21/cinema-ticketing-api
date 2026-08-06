package handler

import (
	scheduledto "cinema-ticketing-api/app/schedule/dto"
	scheduleservice "cinema-ticketing-api/app/schedule/service"
	"cinema-ticketing-api/entities"
	"cinema-ticketing-api/pkg/apperror"
	"cinema-ticketing-api/pkg/pagination"
	"cinema-ticketing-api/pkg/utils"
	"cinema-ticketing-api/request"
	"cinema-ticketing-api/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ScheduleController interface {
	FindAll(c *gin.Context)
	FindUpcoming(c *gin.Context)
	FindUpcomingByTMDBID(c *gin.Context)
	FindByID(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type scheduleController struct {
	scheduleService scheduleservice.ScheduleService
}

func toScheduleResponse(schedule entities.Schedule) scheduledto.ScheduleResponse {
	return scheduledto.ScheduleResponse{
		ID:        schedule.ID,
		MovieID:   schedule.MovieID,
		TMDBID:    schedule.Movie.TMDBID,
		StudioID:  schedule.StudioID,
		Studio:    scheduledto.StudioResponse{ID: schedule.Studio.ID, Name: schedule.Studio.Name},
		StartTime: utils.FormatTime(schedule.StartTime),
		EndTime:   utils.FormatTime(schedule.EndTime),
		Price:     schedule.Price,
	}
}

// Create godoc
// @Summary      Buat jadwal film
// @Description  Membuat jadwal dan menghitung end_time dari durasi film
// @Tags         schedule
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateSchedule true "Schedule data"
// @Success      201 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      404 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /schedule [post]
func (s *scheduleController) Create(c *gin.Context) {
	var req scheduledto.CreateSchedule

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.NewBadRequestError(err.Error()))
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
		TMDBID:    schedule.Movie.TMDBID,
		StudioID:  schedule.StudioID,
		StartTime: utils.FormatTime(schedule.StartTime),
		EndTime:   utils.FormatTime(schedule.EndTime),
		Price:     schedule.Price,
	}))
}

// FindUpcoming godoc
// @Summary      Daftar jadwal mendatang
// @Description  Mengembalikan semua jadwal mendatang beserta TMDB ID dan studio
// @Tags         schedule
// @Produce      json
// @Success      200 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /schedule/upcoming [get]
func (s *scheduleController) FindUpcoming(c *gin.Context) {
	schedules, err := s.scheduleService.FindUpcoming()
	if err != nil {
		c.Error(err)
		return
	}

	result := make([]scheduledto.ScheduleResponse, 0, len(schedules))
	for _, schedule := range schedules {
		result = append(result, toScheduleResponse(schedule))
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Upcoming schedules fetched successfully", result))
}

// FindUpcomingByTMDBID godoc
// @Summary      Jadwal film berdasarkan TMDB ID
// @Description  Mengembalikan jadwal mendatang untuk satu film TMDB
// @Tags         schedule
// @Produce      json
// @Param        tmdb_id path int true "TMDB movie ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /schedule/movie/tmdb/{tmdb_id} [get]
func (s *scheduleController) FindUpcomingByTMDBID(c *gin.Context) {
	tmdbID, err := strconv.ParseInt(c.Param("tmdb_id"), 10, 64)
	if err != nil || tmdbID <= 0 {
		c.Error(apperror.NewBadRequestError("invalid tmdb_id"))
		return
	}

	schedules, err := s.scheduleService.FindUpcomingByTMDBID(tmdbID)
	if err != nil {
		c.Error(err)
		return
	}

	result := make([]scheduledto.ScheduleResponse, 0, len(schedules))
	for _, schedule := range schedules {
		result = append(result, toScheduleResponse(schedule))
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Movie schedules fetched successfully", result))
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
		scheduleResponses = append(scheduleResponses, toScheduleResponse(sch))
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
	c.JSON(http.StatusOK, response.SuccessResponse("Schedule fetched successfully", toScheduleResponse(schedule)))
}

// Update godoc
// @Summary      Perbarui jadwal film
// @Description  Memperbarui jadwal dan menghitung ulang end_time dari durasi film
// @Tags         schedule
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Schedule UUID"
// @Param        request body dto.UpdateSchedule true "Schedule changes"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      404 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /schedule/{id} [put]
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
		c.Error(apperror.NewBadRequestError(err.Error()))
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
		TMDBID:    schedule.Movie.TMDBID,
		StudioID:  schedule.StudioID,
		StartTime: utils.FormatTime(schedule.StartTime),
		EndTime:   utils.FormatTime(schedule.EndTime),
		Price:     schedule.Price,
	}))
}

func NewScheduleController(scheduleService scheduleservice.ScheduleService) ScheduleController {
	return &scheduleController{scheduleService: scheduleService}
}
