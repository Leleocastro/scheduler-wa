package ports

import (
	"complete-api/internal/core/domain"
	"time"
)

type APIGatewayRepository interface {
	CreateConsumer(username, customID string) error
	RateLimitConsumer(username, route string, rateLimit int) error
	CreateACL(username, group string) error
	CreateAPIKey(username string) error
	RemoveRateLimitConsumer(username, route string) error
	RemoveACL(username, group string) error
	GetAPIKey(username string) (string, error)
}

type PaymentRepository interface {
	ValidateSignature(payload []byte, sigHeader string) error
	GetPlanByPriceID(priceID string) (string, error)
	GetEmailByID(customerID string) (string, error)
}

type StatsRepository interface {
	GetUsageByConsumer(username string, startDate, endDate int64) (domain.UsageResponse, error)
}

type RedisRepository interface {
	Get(key string) (string, error)
	Set(key string, value string, expiration time.Duration) error
	Delete(key string) error
	AddScheduledMessage(userID string, schedule domain.ScheduleMessage) error
	GetZRangeByScore(userID string, min, max int64) ([]domain.ScheduleMessage, error)
	RemoveZMember(userID string, member string) error
	GetAllScheduledUsers() ([]string, error)
	UpdateScheduledMessage(userID, id string, msg domain.ScheduleMessage) error
}
