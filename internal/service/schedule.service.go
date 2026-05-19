package service

import (
	"cinema-ticketing-api/internal/model"
	"cinema-ticketing-api/internal/repository"
	"cinema-ticketing-api/internal/request"
	"cinema-ticketing-api/pkg/utils"
	"errors"
	"time"

	"github.com/google/uuid"
)

type ScheduleService interface {
	FindAll(req request.PaginationRequest) ([]model.Schedule, int64, error)
	FindByID(id uuid.UUID) (model.Schedule, error)
	Create(req request.CreateSchedule) (model.Schedule, error)
	Update(id uuid.UUID, req request.UpdateSchedule) (model.Schedule, error)
	Delete(id uuid.UUID) error
}

type scheduleService struct {
	scheduleRepo repository.ScheduleRepository
}

// Create implements [ScheduleService].
func (s *scheduleService) Create(req request.CreateSchedule) (model.Schedule, error) {

	startTime, err := utils.ParseTime(req.StartTime)
	if err != nil {
		return model.Schedule{}, errors.New("invalid start_time format. use 'YYYY-MM-DD HH:MM'")
	}

	endTime, err := utils.ParseTime(req.EndTime)
	if err != nil {
		return model.Schedule{}, errors.New("invalid end_time format. use 'YYYY-MM-DD HH:MM'")
	}

	overlap, err := s.scheduleRepo.CheckScheduleOverlap(req.StudioID, startTime, endTime)
	if err != nil {
		return model.Schedule{}, err
	}
	if overlap {
		return model.Schedule{}, errors.New("schedule overlap")
	}
	return s.scheduleRepo.Create(&model.Schedule{
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
		return errors.New("schedule not found")
	}
	return s.scheduleRepo.Delete(id)
}

// FindAll implements [ScheduleService].
func (s *scheduleService) FindAll(req request.PaginationRequest) ([]model.Schedule, int64, error) {
	return s.scheduleRepo.FindAll(req.Page, req.PerPage)
}

// FindByID implements [ScheduleService].
func (s *scheduleService) FindByID(id uuid.UUID) (model.Schedule, error) {
	return s.scheduleRepo.FindByID(id)
}

// Update implements [ScheduleService].
func (s *scheduleService) Update(id uuid.UUID, req request.UpdateSchedule) (model.Schedule, error) {
	schedule, err := s.scheduleRepo.FindByID(id)
	if err != nil {
		return model.Schedule{}, err
	}
	if schedule.ID == uuid.Nil {
		return model.Schedule{}, errors.New("schedule not found")
	}

	if req.StartTime != "" {
		startTime, err := utils.ParseTime(req.StartTime)
		if err != nil {
			return model.Schedule{}, errors.New("invalid start_time format. use 'YYYY-MM-DD HH:MM'")
		}
		schedule.StartTime = startTime
	}
	if req.EndTime != "" {
		endTime, err := utils.ParseTime(req.EndTime)
		if err != nil {
			return model.Schedule{}, errors.New("invalid end_time format. use 'YYYY-MM-DD HH:MM'")
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

func NewScheduleService(scheduleRepo repository.ScheduleRepository) ScheduleService {
	return &scheduleService{scheduleRepo: scheduleRepo}
}
