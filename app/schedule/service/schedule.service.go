package service

import (
	"errors"
	"math"
	"time"

	scheduledto "cinema-ticketing-api/app/schedule/dto"
	schedulerepository "cinema-ticketing-api/app/schedule/repository"
	"cinema-ticketing-api/entities"
	"cinema-ticketing-api/pkg/apperror"
	"cinema-ticketing-api/pkg/utils"
	"cinema-ticketing-api/request"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScheduleService interface {
	FindAll(req request.PaginationRequest) ([]entities.Schedule, int64, error)
	FindUpcoming() ([]entities.Schedule, error)
	FindUpcomingByTMDBID(tmdbID int64) ([]entities.Schedule, error)
	FindByID(id uuid.UUID) (entities.Schedule, error)
	Create(req scheduledto.CreateSchedule) (entities.Schedule, error)
	Update(id uuid.UUID, req scheduledto.UpdateSchedule) (entities.Schedule, error)
	Delete(id uuid.UUID) error
}

type movieRepository interface {
	FindByID(id uuid.UUID) (*entities.Movie, error)
}

type studioRepository interface {
	FindByID(id uuid.UUID) (*entities.Studio, error)
}

type scheduleService struct {
	scheduleRepo schedulerepository.ScheduleRepository
	movieRepo    movieRepository
	studioRepo   studioRepository
}

func (s *scheduleService) Create(req scheduledto.CreateSchedule) (entities.Schedule, error) {
	startTime, err := utils.ParseTime(req.StartTime)
	if err != nil {
		return entities.Schedule{}, apperror.NewBadRequestError("invalid start_time format. use 'YYYY-MM-DD HH:MM'")
	}
	if !startTime.After(time.Now()) {
		return entities.Schedule{}, apperror.NewBadRequestError("start_time must be in the future")
	}
	if err := validatePrice(req.Price); err != nil {
		return entities.Schedule{}, err
	}

	movie, err := s.findMovie(req.MovieID)
	if err != nil {
		return entities.Schedule{}, err
	}
	studio, err := s.findStudio(req.StudioID)
	if err != nil {
		return entities.Schedule{}, err
	}
	endTime := startTime.Add(time.Duration(movie.Duration) * time.Minute)

	overlap, err := s.scheduleRepo.CheckScheduleOverlap(req.StudioID, startTime, endTime)
	if err != nil {
		return entities.Schedule{}, err
	}
	if overlap {
		return entities.Schedule{}, apperror.NewBadRequestError("schedule overlap")
	}

	schedule, err := s.scheduleRepo.Create(&entities.Schedule{
		ID:        uuid.New(),
		MovieID:   req.MovieID,
		StudioID:  req.StudioID,
		StartTime: startTime,
		EndTime:   endTime,
		Price:     req.Price,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return entities.Schedule{}, err
	}
	schedule.Movie = *movie
	schedule.Studio = *studio
	return schedule, nil
}

func (s *scheduleService) Delete(id uuid.UUID) error {
	schedule, err := s.scheduleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NewNotFoundError("schedule not found")
		}
		return err
	}
	if schedule.ID == uuid.Nil {
		return apperror.NewNotFoundError("schedule not found")
	}
	return s.scheduleRepo.Delete(id)
}

func (s *scheduleService) FindAll(req request.PaginationRequest) ([]entities.Schedule, int64, error) {
	return s.scheduleRepo.FindAll(req.Page, req.PerPage)
}

func (s *scheduleService) FindUpcoming() ([]entities.Schedule, error) {
	return s.scheduleRepo.FindUpcoming(time.Now())
}

func (s *scheduleService) FindUpcomingByTMDBID(tmdbID int64) ([]entities.Schedule, error) {
	if tmdbID <= 0 {
		return nil, apperror.NewBadRequestError("tmdb_id must be greater than zero")
	}
	return s.scheduleRepo.FindUpcomingByTMDBID(tmdbID, time.Now())
}

func (s *scheduleService) FindByID(id uuid.UUID) (entities.Schedule, error) {
	schedule, err := s.scheduleRepo.FindByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entities.Schedule{}, apperror.NewNotFoundError("schedule not found")
	}
	return schedule, err
}

func (s *scheduleService) Update(id uuid.UUID, req scheduledto.UpdateSchedule) (entities.Schedule, error) {
	schedule, err := s.scheduleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Schedule{}, apperror.NewNotFoundError("schedule not found")
		}
		return entities.Schedule{}, err
	}

	if req.MovieID != nil {
		schedule.MovieID = *req.MovieID
	}
	if req.StudioID != nil {
		schedule.StudioID = *req.StudioID
	}
	if req.StartTime != nil {
		startTime, err := utils.ParseTime(*req.StartTime)
		if err != nil {
			return entities.Schedule{}, apperror.NewBadRequestError("invalid start_time format. use 'YYYY-MM-DD HH:MM'")
		}
		schedule.StartTime = startTime
	}
	if req.Price != nil {
		schedule.Price = *req.Price
	}
	if !schedule.StartTime.After(time.Now()) {
		return entities.Schedule{}, apperror.NewBadRequestError("start_time must be in the future")
	}
	if err := validatePrice(schedule.Price); err != nil {
		return entities.Schedule{}, err
	}

	movie, err := s.findMovie(schedule.MovieID)
	if err != nil {
		return entities.Schedule{}, err
	}
	studio, err := s.findStudio(schedule.StudioID)
	if err != nil {
		return entities.Schedule{}, err
	}
	schedule.EndTime = schedule.StartTime.Add(time.Duration(movie.Duration) * time.Minute)

	overlap, err := s.scheduleRepo.CheckScheduleOverlapExcluding(schedule.ID, schedule.StudioID, schedule.StartTime, schedule.EndTime)
	if err != nil {
		return entities.Schedule{}, err
	}
	if overlap {
		return entities.Schedule{}, apperror.NewBadRequestError("schedule overlap")
	}
	schedule.UpdatedAt = time.Now()

	updated, err := s.scheduleRepo.Update(&schedule)
	if err != nil {
		return entities.Schedule{}, err
	}
	updated.Movie = *movie
	updated.Studio = *studio
	return updated, nil
}

func (s *scheduleService) findMovie(id uuid.UUID) (*entities.Movie, error) {
	movie, err := s.movieRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFoundError("movie not found")
		}
		return nil, err
	}
	if movie.Duration <= 0 {
		return nil, apperror.NewBadRequestError("movie duration must be greater than zero")
	}
	return movie, nil
}

func (s *scheduleService) findStudio(id uuid.UUID) (*entities.Studio, error) {
	studio, err := s.studioRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFoundError("studio not found")
		}
		return nil, err
	}
	return studio, nil
}

func validatePrice(price float64) error {
	if price <= 0 || math.Trunc(price) != price {
		return apperror.NewBadRequestError("price must be a positive whole IDR value")
	}
	return nil
}

func NewScheduleService(
	scheduleRepo schedulerepository.ScheduleRepository,
	movieRepo movieRepository,
	studioRepo studioRepository,
) ScheduleService {
	return &scheduleService{scheduleRepo: scheduleRepo, movieRepo: movieRepo, studioRepo: studioRepo}
}
