package job

import (
	schedulerepository "cinema-ticketing-api/app/schedule/repository"
	ticketrepository "cinema-ticketing-api/app/ticket/repository"
	ticketservice "cinema-ticketing-api/app/ticket/service"
	"cinema-ticketing-api/pkg/mailer"
	"log"
	"time"
)

type Scheduler struct {
	transactionService ticketservice.TransactionService
	scheduleRepo       schedulerepository.ScheduleRepository
	ticketRepo         ticketrepository.TicketRepository
	mailer             *mailer.Mailer
	paymentExpiry      time.Duration
}

func NewScheduler(
	transactionService ticketservice.TransactionService,
	scheduleRepo schedulerepository.ScheduleRepository,
	ticketRepo ticketrepository.TicketRepository,
	mailer *mailer.Mailer,
	paymentExpiry time.Duration,
) *Scheduler {
	return &Scheduler{
		transactionService: transactionService,
		scheduleRepo:       scheduleRepo,
		ticketRepo:         ticketRepo,
		mailer:             mailer,
		paymentExpiry:      paymentExpiry,
	}
}

// Start menjalankan semua background job dalam goroutine.
func (s *Scheduler) Start() {
	go s.runAutoCancelExpiredTransactions()
	go s.runUpcomingFilmReminder()
	log.Printf("Scheduler started: auto-cancel (%s) & film reminder (30 min before)", s.paymentExpiry)
}

// runAutoCancelExpiredTransactions menjalankan pengecekan setiap menit.
// Transaksi yang sudah lebih dari 15 menit masih 'pending' akan di-cancel otomatis.
func (s *Scheduler) runAutoCancelExpiredTransactions() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.transactionService.AutoCancelExpiredTransactions()
	}
}

// runUpcomingFilmReminder menjalankan pengecekan setiap menit.
// User yang punya tiket aktif 30 menit sebelum film mulai akan dapat notifikasi.
func (s *Scheduler) runUpcomingFilmReminder() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.sendUpcomingFilmReminders()
	}
}

// sendUpcomingFilmReminders mencari schedule yang akan mulai dalam rentang
// 29-31 menit ke depan, lalu mengirim email reminder ke setiap user yang
// memiliki tiket aktif (paid) pada schedule tersebut.
func (s *Scheduler) sendUpcomingFilmReminders() {
	now := time.Now()
	reminderStart := now.Add(29 * time.Minute)
	reminderEnd := now.Add(31 * time.Minute)

	// Cari schedule yang akan mulai dalam rentang 29-31 menit ke depan
	schedules, err := s.scheduleRepo.FindUpcomingSchedules(reminderStart, reminderEnd)
	if err != nil {
		log.Println("[Scheduler] Error fetching upcoming schedules:", err)
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
			user := ticket.User
			movieTitle := schedule.Movie.Title
			studioName := schedule.Studio.Name
			go func() {
				if err := s.mailer.SendUpcomingFilmReminder(
					user.Email, user.Name,
					movieTitle,
					scheduleTime,
					studioName,
				); err != nil {
					log.Printf("[Scheduler] Failed to send reminder to %s: %v\n", user.Email, err)
				} else {
					log.Printf("[Scheduler] Sent film reminder to %s for %s\n", user.Email, movieTitle)
				}
			}()
		}

		// Log progress
		log.Printf("[Scheduler] Processed reminders for schedule %s: %s at %s\n",
			schedule.ID, schedule.Movie.Title, schedule.StartTime.Format(time.RFC3339))
	}
}
