package ports

import "complete-api/internal/core/domain"

type CheckoutService interface {
	Create(event domain.Event, limit int) error
}
