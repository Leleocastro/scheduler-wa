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

	plan, err := getPlan(planName)
	if err != nil {
		return err
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

func (s *service) CancelSubscription(subscription domain.SubscriptionRoot, username string) error {

	nickname, ok := subscription.Data.Object.Plan.Metadata["nickname"].(string)
	if !ok {
		return fmt.Errorf("nickname is not a string")
	}
	plan, err := getPlan(nickname)
	if err != nil {
		return err
	}

	fmt.Println("Removendo ACL...")

	if err := s.kongRepo.RemoveACL(username, plan.Group); err != nil {
		return err
	}

	fmt.Println("Removendo Rate Limit...")

	if err := s.kongRepo.RemoveRateLimitConsumer(username, plan.Route); err != nil {
		return err
	}

	return nil
}

func getPlan(planName string) (domain.Plan, error) {

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
		return domain.Plan{}, fmt.Errorf("invalid plan name")
	}

	return plan, nil
}

type PlanName string

const (
	Basic    PlanName = "Basic"
	Premium  PlanName = "Premium"
	Business PlanName = "Business"
)
