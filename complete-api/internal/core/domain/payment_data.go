package domain

type Event struct {
	ID              string  `json:"id"`
	Object          string  `json:"object"`
	APIVersion      string  `json:"api_version"`
	Created         int64   `json:"created"`
	Data            Data    `json:"data"`
	Livemode        bool    `json:"livemode"`
	PendingWebhooks int     `json:"pending_webhooks"`
	Request         Request `json:"request"`
	Type            string  `json:"type"`
}

type Data struct {
	Object CheckoutSession `json:"object"`
}

type CheckoutSession struct {
	ID                                string                            `json:"id"`
	Object                            string                            `json:"object"`
	AdaptivePricing                   interface{}                       `json:"adaptive_pricing"`
	AfterExpiration                   interface{}                       `json:"after_expiration"`
	AllowPromotionCodes               interface{}                       `json:"allow_promotion_codes"`
	AmountSubtotal                    int                               `json:"amount_subtotal"`
	AmountTotal                       int                               `json:"amount_total"`
	AutomaticTax                      AutomaticTax                      `json:"automatic_tax"`
	BillingAddressCollection          interface{}                       `json:"billing_address_collection"`
	CancelURL                         string                            `json:"cancel_url"`
	ClientReferenceID                 interface{}                       `json:"client_reference_id"`
	ClientSecret                      interface{}                       `json:"client_secret"`
	Consent                           interface{}                       `json:"consent"`
	ConsentCollection                 interface{}                       `json:"consent_collection"`
	Created                           int64                             `json:"created"`
	Currency                          string                            `json:"currency"`
	CurrencyConversion                interface{}                       `json:"currency_conversion"`
	CustomFields                      []interface{}                     `json:"custom_fields"`
	CustomText                        CustomText                        `json:"custom_text"`
	Customer                          string                            `json:"customer"`
	CustomerCreation                  string                            `json:"customer_creation"`
	CustomerDetails                   CustomerDetails                   `json:"customer_details"`
	CustomerEmail                     string                            `json:"customer_email"`
	ExpiresAt                         int64                             `json:"expires_at"`
	Invoice                           string                            `json:"invoice"`
	InvoiceCreation                   interface{}                       `json:"invoice_creation"`
	Livemode                          bool                              `json:"livemode"`
	Locale                            interface{}                       `json:"locale"`
	Metadata                          Metadata                          `json:"metadata"`
	Mode                              string                            `json:"mode"`
	PaymentIntent                     interface{}                       `json:"payment_intent"`
	PaymentLink                       interface{}                       `json:"payment_link"`
	PaymentMethodCollection           string                            `json:"payment_method_collection"`
	PaymentMethodConfigurationDetails PaymentMethodConfigurationDetails `json:"payment_method_configuration_details"`
	PaymentMethodOptions              PaymentMethodOptions              `json:"payment_method_options"`
	PaymentMethodTypes                []string                          `json:"payment_method_types"`
	PaymentStatus                     string                            `json:"payment_status"`
	PhoneNumberCollection             PhoneNumberCollection             `json:"phone_number_collection"`
	RecoveredFrom                     interface{}                       `json:"recovered_from"`
	SavedPaymentMethodOptions         SavedPaymentMethodOptions         `json:"saved_payment_method_options"`
	SetupIntent                       interface{}                       `json:"setup_intent"`
	ShippingAddressCollection         interface{}                       `json:"shipping_address_collection"`
	ShippingCost                      interface{}                       `json:"shipping_cost"`
	ShippingDetails                   interface{}                       `json:"shipping_details"`
	ShippingOptions                   []interface{}                     `json:"shipping_options"`
	Status                            string                            `json:"status"`
	SubmitType                        interface{}                       `json:"submit_type"`
	Subscription                      string                            `json:"subscription"`
	SuccessURL                        string                            `json:"success_url"`
	TotalDetails                      TotalDetails                      `json:"total_details"`
	UIMode                            string                            `json:"ui_mode"`
	URL                               interface{}                       `json:"url"`
}

type AutomaticTax struct {
	Enabled   bool        `json:"enabled"`
	Liability interface{} `json:"liability"`
	Status    interface{} `json:"status"`
}

type CustomText struct {
	AfterSubmit              interface{} `json:"after_submit"`
	ShippingAddress          interface{} `json:"shipping_address"`
	Submit                   interface{} `json:"submit"`
	TermsOfServiceAcceptance interface{} `json:"terms_of_service_acceptance"`
}

type CustomerDetails struct {
	Address   Address       `json:"address"`
	Email     string        `json:"email"`
	Name      string        `json:"name"`
	Phone     interface{}   `json:"phone"`
	TaxExempt string        `json:"tax_exempt"`
	TaxIds    []interface{} `json:"tax_ids"`
}

type Address struct {
	City       interface{} `json:"city"`
	Country    string      `json:"country"`
	Line1      interface{} `json:"line1"`
	Line2      interface{} `json:"line2"`
	PostalCode interface{} `json:"postal_code"`
	State      interface{} `json:"state"`
}

type Metadata struct {
	PriceID string `json:"price_id"`
}

type PaymentMethodConfigurationDetails struct {
	ID     string      `json:"id"`
	Parent interface{} `json:"parent"`
}

type PaymentMethodOptions struct {
	Card CardOptions `json:"card"`
}

type CardOptions struct {
	RequestThreeDSecure string `json:"request_three_d_secure"`
}

type PhoneNumberCollection struct {
	Enabled bool `json:"enabled"`
}

type SavedPaymentMethodOptions struct {
	AllowRedisplayFilters []string    `json:"allow_redisplay_filters"`
	PaymentMethodRemove   interface{} `json:"payment_method_remove"`
	PaymentMethodSave     interface{} `json:"payment_method_save"`
}

type TotalDetails struct {
	AmountDiscount int `json:"amount_discount"`
	AmountShipping int `json:"amount_shipping"`
	AmountTax      int `json:"amount_tax"`
}

type Request struct {
	ID             interface{} `json:"id"`
	IdempotencyKey interface{} `json:"idempotency_key"`
}
