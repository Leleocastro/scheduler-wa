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

func (s *service) GetAPIKey(username string) (string, error) {
	fmt.Println("Buscando chave de API no Kong...")

	apiKey, err := s.kongRepo.GetAPIKey(username)
	if err != nil {
		return "", err
	}

	return apiKey, nil
}
