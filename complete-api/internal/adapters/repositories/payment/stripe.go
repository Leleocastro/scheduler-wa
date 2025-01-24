package payment

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	stripe "github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/customer"
	"github.com/stripe/stripe-go/v81/price"
	"github.com/stripe/stripe-go/v81/webhook"
)

type StripeClient struct{}

func New(secretKey string) *StripeClient {
	stripe.Key = secretKey

	return &StripeClient{}
}

func (sc *StripeClient) ValidateSignature(payload []byte, sigHeader string) error {
	var secret string

	var webhookPayload WebhookPayload

	if err := json.Unmarshal(payload, &webhookPayload); err != nil {
		return fmt.Errorf("erro ao fazer unmarshal do payload: %v", err)
	}

	whType := webhookPayload.Type

	switch WebhookType(whType) {
	case Checkout:
		secret = os.Getenv("STRIPE_WH_CHECKOUT")
	case Cancel:
		secret = os.Getenv("STRIPE_WH_CANCEL")
	default:
		return nil
	}

	err := webhook.ValidatePayload(payload, sigHeader, secret)
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

func (sc *StripeClient) GetEmailByID(customerID string) (string, error) {
	customer, err := customer.Get(customerID, nil)
	if err != nil {
		log.Printf("Erro ao buscar o cliente: %v", err)
		return "", err
	}

	return customer.Email, nil
}

type WebhookType string

const (
	Checkout WebhookType = "checkout.session.completed"
	Cancel   WebhookType = "customer.subscription.deleted"
)

type WebhookPayload struct {
	Type string `json:"type"`
}
