package service

import (
	"testing"
	"time"

	scheduledto "cinema-ticketing-api/app/schedule/dto"
	"cinema-ticketing-api/entities"
	"cinema-ticketing-api/pkg/utils"

	"github.com/google/uuid"
)

type scheduleRepositoryStub struct {
	created        *entities.Schedule
	requestedTMDB  int64
	requestedAfter time.Time
}

func (r *scheduleRepositoryStub) FindAll(int, int) ([]entities.Schedule, int64, error) {
	return nil, 0, nil
}

func (r *scheduleRepositoryStub) FindUpcomingSchedules(time.Time, time.Time) ([]entities.Schedule, error) {
	return nil, nil
}

func (r *scheduleRepositoryStub) FindUpcoming(time.Time) ([]entities.Schedule, error) {
	return nil, nil
}

func (r *scheduleRepositoryStub) FindUpcomingByTMDBID(tmdbID int64, after time.Time) ([]entities.Schedule, error) {
	r.requestedTMDB = tmdbID
	r.requestedAfter = after
	return []entities.Schedule{{ID: uuid.New()}}, nil
}

func TestFindUpcomingByTMDBIDUsesExternalMovieID(t *testing.T) {
	repository := &scheduleRepositoryStub{}
	service := NewScheduleService(
		repository,
		&scheduleMovieRepositoryStub{},
		&scheduleStudioRepositoryStub{},
	)

	schedules, err := service.FindUpcomingByTMDBID(299534)
	if err != nil {
		t.Fatalf("unexpected find error: %v", err)
	}
	if len(schedules) != 1 || repository.requestedTMDB != 299534 {
		t.Fatalf("expected lookup by TMDB ID 299534")
	}
	if repository.requestedAfter.IsZero() {
		t.Fatal("expected upcoming lookup to include the current time")
	}
}

func (r *scheduleRepositoryStub) FindByID(uuid.UUID) (entities.Schedule, error) {
	return entities.Schedule{}, nil
}

func (r *scheduleRepositoryStub) Create(schedule *entities.Schedule) (entities.Schedule, error) {
	r.created = schedule
	return *schedule, nil
}

func (r *scheduleRepositoryStub) Update(schedule *entities.Schedule) (entities.Schedule, error) {
	return *schedule, nil
}

func (r *scheduleRepositoryStub) Delete(uuid.UUID) error { return nil }

func (r *scheduleRepositoryStub) CheckScheduleOverlap(uuid.UUID, time.Time, time.Time) (bool, error) {
	return false, nil
}

func (r *scheduleRepositoryStub) CheckScheduleOverlapExcluding(uuid.UUID, uuid.UUID, time.Time, time.Time) (bool, error) {
	return false, nil
}

type scheduleMovieRepositoryStub struct {
	movie *entities.Movie
}

func (r *scheduleMovieRepositoryStub) FindByID(uuid.UUID) (*entities.Movie, error) {
	return r.movie, nil
}

type scheduleStudioRepositoryStub struct {
	studio *entities.Studio
}

func (r *scheduleStudioRepositoryStub) FindByID(uuid.UUID) (*entities.Studio, error) {
	return r.studio, nil
}

func TestCreateScheduleDerivesEndTimeFromMovieDuration(t *testing.T) {
	movieID := uuid.New()
	studioID := uuid.New()
	tmdbID := int64(299534)
	repository := &scheduleRepositoryStub{}
	service := NewScheduleService(
		repository,
		&scheduleMovieRepositoryStub{movie: &entities.Movie{ID: movieID, TMDBID: &tmdbID, Duration: 181}},
		&scheduleStudioRepositoryStub{studio: &entities.Studio{ID: studioID, Name: "Studio 1"}},
	)
	start := time.Now().Add(24 * time.Hour)

	schedule, err := service.Create(scheduledto.CreateSchedule{
		MovieID:   movieID,
		StudioID:  studioID,
		StartTime: utils.FormatTime(start),
		Price:     50000,
	})
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}
	if schedule.EndTime.Sub(schedule.StartTime) != 181*time.Minute {
		t.Fatalf("expected 181 minute duration, got %s", schedule.EndTime.Sub(schedule.StartTime))
	}
	if repository.created == nil {
		t.Fatal("expected schedule to be persisted")
	}
}
