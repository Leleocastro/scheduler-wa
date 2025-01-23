package checkoutsrv

import (
	"complete-api/internal/core/domain"
	"complete-api/internal/core/ports"
	"fmt"
)

type service struct {
	kongRepo ports.APIGatewayRepository
}

func New(kongRepo ports.APIGatewayRepository) *service {
	return &service{
		kongRepo: kongRepo,
	}
}

func (s *service) Create(event domain.Event, planName string) error {
	username := event.Data.Object.CustomerEmail

	var plan domain.Plan

	switch PlanName(planName) {
	case Basic:
		plan = domain.Plan{
			Name:        "Basic",
			WebSocket:   false,
			LimitPerDay: 100,
			Route:       "whatsapp-route",
			Group:       "basic",
		}
	case Premium:
		plan = domain.Plan{
			Name:        "Premium",
			WebSocket:   true,
			LimitPerDay: 100,
			Route:       "whatsapp-route",
			Group:       "premium",
		}
	case Business:
		plan = domain.Plan{
			Name:        "Business",
			WebSocket:   true,
			LimitPerDay: 100000,
			Route:       "whatsapp-route",
			Group:       "business",
		}
	default:
		return fmt.Errorf("invalid plan name")
	}

	fmt.Println("Adicionando Rate Limit...")

	if err := s.kongRepo.RateLimitConsumer(username, plan.Route, plan.LimitPerDay); err != nil {
		return err
	}

	if err := s.kongRepo.CreateACL(username, plan.Group); err != nil {
		return err
	}

	if err := s.kongRepo.CreateAPIKey(username); err != nil {
		return err
	}

	return nil
}

type PlanName string

const (
	Basic    PlanName = "Basic"
	Premium  PlanName = "Premium"
	Business PlanName = "Business"
)
