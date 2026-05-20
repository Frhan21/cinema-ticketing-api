package scheduler

import (
	"cinema-ticketing-api/internal/enums"
	"cinema-ticketing-api/internal/model"
	"cinema-ticketing-api/pkg/mailer"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type Scheduler struct {
	db     *gorm.DB
	mailer *mailer.Mailer
}

func NewScheduler(db *gorm.DB, mailer *mailer.Mailer) *Scheduler {
	return &Scheduler{db: db, mailer: mailer}
}

// Start menjalankan semua background job dalam goroutine.
func (s *Scheduler) Start() {
	go s.runAutoCancelExpiredTransactions()
	go s.runUpcomingFilmReminder()
	log.Println("Scheduler started: auto-cancel (15 min) & film reminder (30 min before)")
}

// runAutoCancelExpiredTransactions menjalankan pengecekan setiap menit.
// Transaksi yang sudah lebih dari 15 menit masih 'pending' akan di-cancel otomatis.
func (s *Scheduler) runAutoCancelExpiredTransactions() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.autoCancelExpiredTransactions()
	}
}

func (s *Scheduler) autoCancelExpiredTransactions() {
	expiredBefore := time.Now().Add(-15 * time.Minute)

	var expiredTransactions []model.Transaction
	if err := s.db.Preload("Items").Preload("User").
		Where("payment_status = ? AND created_at < ?", enums.PaymentStatusPending, expiredBefore).
		Find(&expiredTransactions).Error; err != nil {
		log.Println("[Scheduler] Error fetching expired transactions:", err)
		return
	}

	for _, tx := range expiredTransactions {
		// Update status transaksi menjadi cancelled
		if err := s.db.Model(&model.Transaction{}).
			Where("id = ?", tx.ID).
			Update("payment_status", enums.PaymentStatusCancelled).Error; err != nil {
			log.Printf("[Scheduler] Failed to cancel transaction %s: %v\n", tx.ID, err)
			continue
		}

		// Update semua tiket terkait menjadi cancelled
		var ticketIDs []interface{}
		for _, item := range tx.Items {
			ticketIDs = append(ticketIDs, item.TicketID)
		}
		if len(ticketIDs) > 0 {
			s.db.Model(&model.Ticket{}).
				Where("id IN ?", ticketIDs).
				Update("status", enums.TicketStatusCancelled)
		}

		log.Printf("[Scheduler] Auto-cancelled transaction %s (expired 15 min)\n", tx.ID)

		// Kirim notifikasi email ke user (non-blocking)
		go func(user model.User) {
			if err := s.mailer.SendTicketCancelledNotification(
				user.Email, user.Name, "Tiket Anda",
			); err != nil {
				log.Printf("[Scheduler] Failed to send cancel email to %s: %v\n", user.Email, err)
			}
		}(tx.User)
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

func (s *Scheduler) sendUpcomingFilmReminders() {
	now := time.Now()
	reminderStart := now.Add(29 * time.Minute)
	reminderEnd := now.Add(31 * time.Minute)

	// Cari schedule yang akan mulai dalam rentang 29-31 menit ke depan
	var schedules []model.Schedule
	if err := s.db.Preload("Movie").Preload("Studio").
		Where("start_time BETWEEN ? AND ?", reminderStart, reminderEnd).
		Find(&schedules).Error; err != nil {
		log.Println("[Scheduler] Error fetching upcoming schedules:", err)
		return
	}

	for _, schedule := range schedules {
		// Cari semua tiket aktif (paid) untuk schedule ini beserta user-nya
		var tickets []model.Ticket
		if err := s.db.Preload("User").
			Where("schedule_id = ? AND status = ?", schedule.ID, enums.TicketStatusPaid).
			Find(&tickets).Error; err != nil {
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
					log.Printf("[Scheduler] Failed to send reminder to %s: %v\n", user.Email, err)
				} else {
					log.Printf("[Scheduler] Sent film reminder to %s for %s\n", user.Email, schedule.Movie.Title)
				}
			}(ticket.User)
		}

		// Log progress
		log.Printf("[Scheduler] Processed reminders for schedule %s: %s at %s\n",
			schedule.ID, schedule.Movie.Title, fmt.Sprint(schedule.StartTime))
	}
}
