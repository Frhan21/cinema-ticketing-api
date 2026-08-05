package service

import (
	scheduledto "cinema-ticketing-api/app/schedule/dto"
	schedulerepository "cinema-ticketing-api/app/schedule/repository"
	"cinema-ticketing-api/entities"
	"cinema-ticketing-api/pkg/apperror"
	"cinema-ticketing-api/pkg/utils"
	"cinema-ticketing-api/request"
	"time"

	"github.com/google/uuid"
)

type ScheduleService interface {
	FindAll(req request.PaginationRequest) ([]entities.Schedule, int64, error)
	FindByID(id uuid.UUID) (entities.Schedule, error)
	Create(req scheduledto.CreateSchedule) (entities.Schedule, error)
	Update(id uuid.UUID, req scheduledto.UpdateSchedule) (entities.Schedule, error)
	Delete(id uuid.UUID) error
}

type scheduleService struct {
	scheduleRepo schedulerepository.ScheduleRepository
}

// Create implements [ScheduleService].
func (s *scheduleService) Create(req scheduledto.CreateSchedule) (entities.Schedule, error) {

	startTime, err := utils.ParseTime(req.StartTime)
	if err != nil {
		return entities.Schedule{}, apperror.NewBadRequestError("invalid start_time format. use 'YYYY-MM-DD HH:MM'")
	}

	endTime, err := utils.ParseTime(req.EndTime)
	if err != nil {
		return entities.Schedule{}, apperror.NewBadRequestError("invalid end_time format. use 'YYYY-MM-DD HH:MM'")
	}

	overlap, err := s.scheduleRepo.CheckScheduleOverlap(req.StudioID, startTime, endTime)
	if err != nil {
		return entities.Schedule{}, err
	}
	if overlap {
		return entities.Schedule{}, apperror.NewBadRequestError("schedule overlap")
	}
	return s.scheduleRepo.Create(&entities.Schedule{
		ID:        uuid.New(),
		MovieID:   req.MovieID,
		StudioID:  req.StudioID,
		StartTime: startTime,
		EndTime:   endTime,
		Price:     req.Price,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}

// Delete implements [ScheduleService].
func (s *scheduleService) Delete(id uuid.UUID) error {
	schedule, err := s.scheduleRepo.FindByID(id)
	if err != nil {
		return err
	}
	if schedule.ID == uuid.Nil {
		return apperror.NewNotFoundError("schedule not found")
	}
	return s.scheduleRepo.Delete(id)
}

// FindAll implements [ScheduleService].
func (s *scheduleService) FindAll(req request.PaginationRequest) ([]entities.Schedule, int64, error) {
	return s.scheduleRepo.FindAll(req.Page, req.PerPage)
}

// FindByID implements [ScheduleService].
func (s *scheduleService) FindByID(id uuid.UUID) (entities.Schedule, error) {
	return s.scheduleRepo.FindByID(id)
}

// Update implements [ScheduleService].
func (s *scheduleService) Update(id uuid.UUID, req scheduledto.UpdateSchedule) (entities.Schedule, error) {
	schedule, err := s.scheduleRepo.FindByID(id)
	if err != nil {
		return entities.Schedule{}, err
	}
	if schedule.ID == uuid.Nil {
		return entities.Schedule{}, apperror.NewNotFoundError("schedule not found")
	}

	if req.StartTime != "" {
		startTime, err := utils.ParseTime(req.StartTime)
		if err != nil {
			return entities.Schedule{}, apperror.NewBadRequestError("invalid start_time format. use 'YYYY-MM-DD HH:MM'")
		}
		schedule.StartTime = startTime
	}
	if req.EndTime != "" {
		endTime, err := utils.ParseTime(req.EndTime)
		if err != nil {
			return entities.Schedule{}, apperror.NewBadRequestError("invalid end_time format. use 'YYYY-MM-DD HH:MM'")
		}
		schedule.EndTime = endTime
	}
	if req.Price != 0 {
		schedule.Price = req.Price
	}
	if req.MovieID != uuid.Nil {
		schedule.MovieID = req.MovieID
	}
	if req.StudioID != uuid.Nil {
		schedule.StudioID = req.StudioID
	}
	schedule.UpdatedAt = time.Now()
	return s.scheduleRepo.Update(&schedule)
}

func NewScheduleService(scheduleRepo schedulerepository.ScheduleRepository) ScheduleService {
	return &scheduleService{scheduleRepo: scheduleRepo}
}
