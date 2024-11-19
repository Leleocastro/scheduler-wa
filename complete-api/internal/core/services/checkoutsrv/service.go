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

func (s *service) Create(event domain.Event, limit int) error {
	username := event.Data.Object.Metadata["username"]

	//TODO: Buscar o plano e com base no plano planejar qual a rota que será liberada

	fmt.Println("Criando consumidor no Kong...")

	if err := s.kongRepo.RateLimitConsumer(username, "whatsapp", limit); err != nil {
		return err
	}

	return nil
}
