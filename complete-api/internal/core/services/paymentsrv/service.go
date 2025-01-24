package paymentsrv

import "complete-api/internal/core/ports"

type service struct {
	paymentRepo ports.PaymentRepository
}

func New(paymentRepo ports.PaymentRepository) *service {
	return &service{
		paymentRepo: paymentRepo,
	}
}

func (s *service) ValidateSignature(payload []byte, sigHeader string) error {
	return s.paymentRepo.ValidateSignature(payload, sigHeader)
}

func (s *service) GetPlanByPriceID(priceID string) (string, error) {
	return s.paymentRepo.GetPlanByPriceID(priceID)
}

func (s *service) GetEmailByID(customerID string) (string, error) {
	return s.paymentRepo.GetEmailByID(customerID)
}
