package repository

import (
	"cinema-ticketing-api/internal/model"
	"cinema-ticketing-api/pkg/pagination"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScheduleRepository interface {
	FindAll(page, perPage int) ([]model.Schedule, int64, error)
	FindByID(id uuid.UUID) (model.Schedule, error)
	Create(schedule *model.Schedule) (model.Schedule, error)
	Update(schedule *model.Schedule) (model.Schedule, error)
	Delete(id uuid.UUID) error
	CheckScheduleOverlap(studioID uuid.UUID, startTime, endTime time.Time) (bool, error)
}

type scheduleRepository struct {
	DB *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) ScheduleRepository {
	return &scheduleRepository{DB: db}
}

func (r *scheduleRepository) FindAll(page, perPage int) ([]model.Schedule, int64, error) {
	var schedules []model.Schedule
	var count int64

	err := r.DB.Model(&model.Schedule{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.DB.Preload("Movie").Preload("Studio").Scopes(pagination.Paginate(page, perPage)).Find(&schedules).Error
	return schedules, count, err
}

func (r *scheduleRepository) FindByID(id uuid.UUID) (model.Schedule, error) {
	var schedule model.Schedule
	err := r.DB.Preload("Movie").Preload("Studio").Where("id = ?", id).First(&schedule).Error
	return schedule, err
}

func (r *scheduleRepository) Create(schedule *model.Schedule) (model.Schedule, error) {
	err := r.DB.Create(schedule).Error
	return *schedule, err
}

func (r *scheduleRepository) Update(schedule *model.Schedule) (model.Schedule, error) {
	err := r.DB.Save(schedule).Error
	return *schedule, err
}

func (r *scheduleRepository) Delete(id uuid.UUID) error {
	err := r.DB.Delete(&model.Schedule{}, "id = ?", id).Error
	return err
}

func (r *scheduleRepository) CheckScheduleOverlap(studioID uuid.UUID, startTime, endTime time.Time) (bool, error) {
	var count int64
	err := r.DB.Model(&model.Schedule{}).Where("studio_id = ? AND start_time < ? AND end_time > ?", studioID, endTime, startTime).Count(&count).Error
	return count > 0, err
}
