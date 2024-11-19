package domain

type Event struct {
	ID              string    `json:"id"`
	Object          string    `json:"object"`
	APIVersion      string    `json:"api_version"`
	Created         int64     `json:"created"`
	Data            EventData `json:"data"`
	Livemode        bool      `json:"livemode"`
	PendingWebhooks int       `json:"pending_webhooks"`
	Request         Request   `json:"request"`
	Type            string    `json:"type"`
}

type EventData struct {
	Object PaymentIntent `json:"object"`
}

type PaymentIntent struct {
	ID                                string               `json:"id"`
	Object                            string               `json:"object"`
	Amount                            int                  `json:"amount"`
	AmountCapturable                  int                  `json:"amount_capturable"`
	AmountDetails                     AmountDetails        `json:"amount_details"`
	AmountReceived                    int                  `json:"amount_received"`
	Application                       *string              `json:"application"`
	ApplicationFeeAmount              *int                 `json:"application_fee_amount"`
	AutomaticPaymentMethods           *string              `json:"automatic_payment_methods"`
	CanceledAt                        *int64               `json:"canceled_at"`
	CancellationReason                *string              `json:"cancellation_reason"`
	CaptureMethod                     string               `json:"capture_method"`
	ClientSecret                      string               `json:"client_secret"`
	ConfirmationMethod                string               `json:"confirmation_method"`
	Created                           int64                `json:"created"`
	Currency                          string               `json:"currency"`
	Customer                          *string              `json:"customer"`
	Description                       string               `json:"description"`
	Invoice                           *string              `json:"invoice"`
	LastPaymentError                  *string              `json:"last_payment_error"`
	LatestCharge                      string               `json:"latest_charge"`
	Livemode                          bool                 `json:"livemode"`
	Metadata                          map[string]string    `json:"metadata"`
	NextAction                        *string              `json:"next_action"`
	OnBehalfOf                        *string              `json:"on_behalf_of"`
	PaymentMethod                     string               `json:"payment_method"`
	PaymentMethodConfigurationDetails *string              `json:"payment_method_configuration_details"`
	PaymentMethodOptions              PaymentMethodOptions `json:"payment_method_options"`
	PaymentMethodTypes                []string             `json:"payment_method_types"`
	Processing                        *string              `json:"processing"`
	ReceiptEmail                      *string              `json:"receipt_email"`
	Review                            *string              `json:"review"`
	SetupFutureUsage                  *string              `json:"setup_future_usage"`
	Shipping                          Shipping             `json:"shipping"`
	Source                            *string              `json:"source"`
	StatementDescriptor               *string              `json:"statement_descriptor"`
	StatementDescriptorSuffix         *string              `json:"statement_descriptor_suffix"`
	Status                            string               `json:"status"`
	TransferData                      *string              `json:"transfer_data"`
	TransferGroup                     *string              `json:"transfer_group"`
}

type AmountDetails struct {
	Tip map[string]interface{} `json:"tip"`
}

type PaymentMethodOptions struct {
	Card CardOptions `json:"card"`
}

type CardOptions struct {
	Installments        *string `json:"installments"`
	MandateOptions      *string `json:"mandate_options"`
	Network             *string `json:"network"`
	RequestThreeDSecure string  `json:"request_three_d_secure"`
}

type Shipping struct {
	Address        Address `json:"address"`
	Carrier        *string `json:"carrier"`
	Name           string  `json:"name"`
	Phone          *string `json:"phone"`
	TrackingNumber *string `json:"tracking_number"`
}

type Address struct {
	City       string  `json:"city"`
	Country    string  `json:"country"`
	Line1      string  `json:"line1"`
	Line2      *string `json:"line2"`
	PostalCode string  `json:"postal_code"`
	State      string  `json:"state"`
}

type Request struct {
	ID             string `json:"id"`
	IdempotencyKey string `json:"idempotency_key"`
}
