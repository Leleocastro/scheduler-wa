package ports

import "complete-api/internal/core/domain"

type CheckoutService interface {
	Create(event domain.Event, plan string) error
	CancelSubscription(subscription domain.SubscriptionRoot, username string) error
}

type GatewayService interface {
	CreateConsumer(username, customID string) error
	GetAPIKey(username string) (string, error)
}

type PaymentService interface {
	ValidateSignature(payload []byte, sigHeader string) error
	GetPlanByPriceID(priceID string) (string, error)
	GetEmailByID(customerID string) (string, error)
}

type StatsService interface {
	GetUsageByConsumer(username, startDate, endDate string) (domain.Usage, error)
}

type SchedulesService interface {
	GetSchedules(username string) ([]domain.ScheduleMessage, error)
	CreateSchedule(username string, schedule domain.ScheduleMessage) (domain.ScheduleMessage, error)
	UpdateSchedule(username string, scheduleID string, schedule domain.ScheduleMessage) (domain.ScheduleMessage, error)
	DeleteSchedule(username string, scheduleID string) error
}
