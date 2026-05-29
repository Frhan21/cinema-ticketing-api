package service

import (
	"cinema-ticketing-api/internal/model"
	"cinema-ticketing-api/internal/repository"
	"cinema-ticketing-api/internal/request"
	"cinema-ticketing-api/pkg/apperror"
	"cinema-ticketing-api/pkg/mailer"
	"cinema-ticketing-api/pkg/utils"
	"log"
	"time"
	"fmt"

	"github.com/google/uuid"
)

type ScheduleService interface {
	FindAll(req request.PaginationRequest) ([]model.Schedule, int64, error)
	FindByID(id uuid.UUID) (model.Schedule, error)
	Create(req request.CreateSchedule) (model.Schedule, error)
	Update(id uuid.UUID, req request.UpdateSchedule) (model.Schedule, error)
	Delete(id uuid.UUID) error
	SendUpcomingFilmReminders()
}

type scheduleService struct {
	scheduleRepo repository.ScheduleRepository
	ticketRepo   repository.TicketRepository
	mailer       *mailer.Mailer
}

// Create implements [ScheduleService].
func (s *scheduleService) Create(req request.CreateSchedule) (model.Schedule, error) {

	startTime, err := utils.ParseTime(req.StartTime)
	if err != nil {
		return model.Schedule{}, apperror.NewBadRequestError("invalid start_time format. use 'YYYY-MM-DD HH:MM'")
	}

	endTime, err := utils.ParseTime(req.EndTime)
	if err != nil {
		return model.Schedule{}, apperror.NewBadRequestError("invalid end_time format. use 'YYYY-MM-DD HH:MM'")
	}

	overlap, err := s.scheduleRepo.CheckScheduleOverlap(req.StudioID, startTime, endTime)
	if err != nil {
		return model.Schedule{}, err
	}
	if overlap {
		return model.Schedule{}, apperror.NewBadRequestError("schedule overlap")
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
		return apperror.NewNotFoundError("schedule not found")
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
		return model.Schedule{}, apperror.NewNotFoundError("schedule not found")
	}

	if req.StartTime != "" {
		startTime, err := utils.ParseTime(req.StartTime)
		if err != nil {
			return model.Schedule{}, apperror.NewBadRequestError("invalid start_time format. use 'YYYY-MM-DD HH:MM'")
		}
		schedule.StartTime = startTime
	}
	if req.EndTime != "" {
		endTime, err := utils.ParseTime(req.EndTime)
		if err != nil {
			return model.Schedule{}, apperror.NewBadRequestError("invalid end_time format. use 'YYYY-MM-DD HH:MM'")
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

func NewScheduleService(scheduleRepo repository.ScheduleRepository, ticketRepo repository.TicketRepository, mailer *mailer.Mailer) ScheduleService {
	return &scheduleService{scheduleRepo: scheduleRepo, ticketRepo: ticketRepo, mailer: mailer}
}

// SendUpcomingFilmReminders implements [ScheduleService].
func (s *scheduleService) SendUpcomingFilmReminders() {
	now := time.Now()
	reminderStart := now.Add(29 * time.Minute)
	reminderEnd := now.Add(31 * time.Minute)

	// Cari schedule yang akan mulai dalam rentang 29-31 menit ke depan
	schedules, err := s.scheduleRepo.FindUpcomingSchedules(reminderStart, reminderEnd)
	if err != nil {
		log.Println("[ScheduleService] Error fetching upcoming schedules:", err)
		return
	}

	for _, schedule := range schedules {
		// Cari semua tiket aktif (paid) untuk schedule ini beserta user-nya
		tickets, err := s.ticketRepo.FindPaidTicketsBySchedule(schedule.ID)
		if err != nil {
			continue
		}

		scheduleTime := schedule.StartTime.In(time.FixedZone("WIB", 7*3600)).Format("02 Jan 2006, 15:04 WIB")

		for _, ticket := range tickets {
			go func(user model.User) {
				if err := s.mailer.SendUpcomingFilmReminder(
					user.Email, user.Name,
					schedule.Movie.Title,
					scheduleTime,
					schedule.Studio.Name,
				); err != nil {
					log.Printf("[ScheduleService] Failed to send reminder to %s: %v\n", user.Email, err)
				} else {
					log.Printf("[ScheduleService] Sent film reminder to %s for %s\n", user.Email, schedule.Movie.Title)
				}
			}(ticket.User)
		}

		// Log progress
		log.Printf("[ScheduleService] Processed reminders for schedule %s: %s at %s\n",
			schedule.ID, schedule.Movie.Title, fmt.Sprint(schedule.StartTime))
	}
}
