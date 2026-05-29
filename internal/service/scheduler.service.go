package service

import (
	"log"
	"time"
)

type Scheduler struct {
	transactionService TransactionService
	scheduleService    ScheduleService
}

func NewScheduler(transactionService TransactionService, scheduleService ScheduleService) *Scheduler {
	return &Scheduler{
		transactionService: transactionService,
		scheduleService:    scheduleService,
	}
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
		s.transactionService.AutoCancelExpiredTransactions()
	}
}

// runUpcomingFilmReminder menjalankan pengecekan setiap menit.
// User yang punya tiket aktif 30 menit sebelum film mulai akan dapat notifikasi.
func (s *Scheduler) runUpcomingFilmReminder() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.scheduleService.SendUpcomingFilmReminders()
	}
}
