package ports

import "complete-api/internal/core/domain"

type CheckoutService interface {
	Create(event domain.Event, limit int) error
}

type GatewayService interface {
	CreateConsumer(username, customID string) error
}
