package setup

import (
	"cinema-ticketing-api/internal/service"
)

func InitScheduler(transactionService service.TransactionService, scheduleService service.ScheduleService) {
	scheduler := service.NewScheduler(transactionService, scheduleService)
	scheduler.Start()
}
