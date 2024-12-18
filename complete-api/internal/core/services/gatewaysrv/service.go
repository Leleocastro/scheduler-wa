package gatewaysrv

import (
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

func (s *service) CreateConsumer(username, customID string) error {
	fmt.Println("Criando consumidor no Kong...")

	if err := s.kongRepo.CreateConsumer(username, customID); err != nil {
		return err
	}

	return nil
}
