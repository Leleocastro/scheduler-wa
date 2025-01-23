package payment

import (
	stripe "github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/price"
	"github.com/stripe/stripe-go/v81/webhook"
)

type StripeClient struct {
	webhookSecret string
}

func New(secretKey, webhookSecret string) *StripeClient {
	stripe.Key = secretKey

	return &StripeClient{
		webhookSecret: webhookSecret,
	}
}

func (sc *StripeClient) ValidateSignature(payload []byte, sigHeader string) error {
	err := webhook.ValidatePayload(payload, sigHeader, sc.webhookSecret)
	return err
}

func (sc *StripeClient) GetPlanByPriceID(priceID string) (string, error) {
	price, err := price.Get(priceID, nil)
	if err != nil {
		return "", err
	}

	plan := price.Metadata["nickname"]

	return plan, nil
}
