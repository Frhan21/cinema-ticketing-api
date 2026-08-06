package repository

import (
	"cinema-ticketing-api/entities"
	"cinema-ticketing-api/pkg/pagination"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScheduleRepository interface {
	FindAll(page, perPage int) ([]entities.Schedule, int64, error)
	FindUpcomingSchedules(start, end time.Time) ([]entities.Schedule, error)
	FindUpcoming(now time.Time) ([]entities.Schedule, error)
	FindUpcomingByTMDBID(tmdbID int64, now time.Time) ([]entities.Schedule, error)
	FindByID(id uuid.UUID) (entities.Schedule, error)
	Create(schedule *entities.Schedule) (entities.Schedule, error)
	Update(schedule *entities.Schedule) (entities.Schedule, error)
	Delete(id uuid.UUID) error
	CheckScheduleOverlap(studioID uuid.UUID, startTime, endTime time.Time) (bool, error)
	CheckScheduleOverlapExcluding(id, studioID uuid.UUID, startTime, endTime time.Time) (bool, error)
}

type scheduleRepository struct {
	DB *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) ScheduleRepository {
	return &scheduleRepository{DB: db}
}

func (r *scheduleRepository) FindAll(page, perPage int) ([]entities.Schedule, int64, error) {
	var schedules []entities.Schedule
	var count int64

	err := r.DB.Model(&entities.Schedule{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.DB.Preload("Movie").Preload("Studio").Scopes(pagination.Paginate(page, perPage)).Find(&schedules).Error
	return schedules, count, err
}

func (r *scheduleRepository) FindUpcomingSchedules(start, end time.Time) ([]entities.Schedule, error) {
	var schedules []entities.Schedule
	err := r.DB.Preload("Movie").Preload("Studio").
		Where("start_time BETWEEN ? AND ?", start, end).
		Find(&schedules).Error
	return schedules, err
}

func (r *scheduleRepository) FindUpcoming(now time.Time) ([]entities.Schedule, error) {
	var schedules []entities.Schedule
	err := r.DB.Preload("Movie").Preload("Studio").
		Where("start_time > ?", now).
		Order("start_time ASC").
		Find(&schedules).Error
	return schedules, err
}

func (r *scheduleRepository) FindUpcomingByTMDBID(tmdbID int64, now time.Time) ([]entities.Schedule, error) {
	var schedules []entities.Schedule
	err := r.DB.Preload("Movie").Preload("Studio").
		Joins("JOIN movies ON movies.id = schedules.movie_id").
		Where("movies.tmdb_id = ? AND schedules.start_time > ?", tmdbID, now).
		Order("schedules.start_time ASC").
		Find(&schedules).Error
	return schedules, err
}

func (r *scheduleRepository) FindByID(id uuid.UUID) (entities.Schedule, error) {
	var schedule entities.Schedule
	err := r.DB.Preload("Movie").Preload("Studio").Where("id = ?", id).First(&schedule).Error
	return schedule, err
}

func (r *scheduleRepository) Create(schedule *entities.Schedule) (entities.Schedule, error) {
	err := r.DB.Create(schedule).Error
	return *schedule, err
}

func (r *scheduleRepository) Update(schedule *entities.Schedule) (entities.Schedule, error) {
	err := r.DB.Save(schedule).Error
	return *schedule, err
}

func (r *scheduleRepository) Delete(id uuid.UUID) error {
	err := r.DB.Delete(&entities.Schedule{}, "id = ?", id).Error
	return err
}

func (r *scheduleRepository) CheckScheduleOverlap(studioID uuid.UUID, startTime, endTime time.Time) (bool, error) {
	var count int64
	err := r.DB.Model(&entities.Schedule{}).Where("studio_id = ? AND start_time < ? AND end_time > ?", studioID, endTime, startTime).Count(&count).Error
	return count > 0, err
}

func (r *scheduleRepository) CheckScheduleOverlapExcluding(id, studioID uuid.UUID, startTime, endTime time.Time) (bool, error) {
	var count int64
	err := r.DB.Model(&entities.Schedule{}).
		Where("id <> ? AND studio_id = ? AND start_time < ? AND end_time > ?", id, studioID, endTime, startTime).
		Count(&count).Error
	return count > 0, err
}
