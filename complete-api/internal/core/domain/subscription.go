package domain

type SubscriptionRoot struct {
	ID              string     `json:"id"`
	Object          string     `json:"object"`
	APIVersion      string     `json:"api_version"`
	Created         int64      `json:"created"`
	Data            SubData    `json:"data"`
	Livemode        bool       `json:"livemode"`
	PendingWebhooks int        `json:"pending_webhooks"`
	Request         SubRequest `json:"request"`
	Type            string     `json:"type"`
}

type SubData struct {
	Object SubscriptionObject `json:"object"`
}

type SubscriptionObject struct {
	ID                            string              `json:"id"`
	Object                        string              `json:"object"`
	Application                   interface{}         `json:"application"`
	ApplicationFeePercent         interface{}         `json:"application_fee_percent"`
	AutomaticTax                  SubAutomaticTax     `json:"automatic_tax"`
	BillingCycleAnchor            int64               `json:"billing_cycle_anchor"`
	BillingCycleAnchorConfig      interface{}         `json:"billing_cycle_anchor_config"`
	BillingThresholds             interface{}         `json:"billing_thresholds"`
	CancelAt                      interface{}         `json:"cancel_at"`
	CancelAtPeriodEnd             bool                `json:"cancel_at_period_end"`
	CanceledAt                    int64               `json:"canceled_at"`
	CancellationDetails           CancellationDetails `json:"cancellation_details"`
	CollectionMethod              string              `json:"collection_method"`
	Created                       int64               `json:"created"`
	Currency                      string              `json:"currency"`
	CurrentPeriodEnd              int64               `json:"current_period_end"`
	CurrentPeriodStart            int64               `json:"current_period_start"`
	Customer                      string              `json:"customer"`
	DaysUntilDue                  interface{}         `json:"days_until_due"`
	DefaultPaymentMethod          string              `json:"default_payment_method"`
	DefaultSource                 interface{}         `json:"default_source"`
	DefaultTaxRates               []interface{}       `json:"default_tax_rates"`
	Description                   interface{}         `json:"description"`
	Discount                      interface{}         `json:"discount"`
	Discounts                     []interface{}       `json:"discounts"`
	EndedAt                       int64               `json:"ended_at"`
	InvoiceSettings               InvoiceSettings     `json:"invoice_settings"`
	Items                         Items               `json:"items"`
	LatestInvoice                 string              `json:"latest_invoice"`
	Livemode                      bool                `json:"livemode"`
	Metadata                      map[string]string   `json:"metadata"`
	NextPendingInvoiceItemInvoice interface{}         `json:"next_pending_invoice_item_invoice"`
	OnBehalfOf                    interface{}         `json:"on_behalf_of"`
	PauseCollection               interface{}         `json:"pause_collection"`
	PaymentSettings               PaymentSettings     `json:"payment_settings"`
	PendingInvoiceItemInterval    interface{}         `json:"pending_invoice_item_interval"`
	PendingSetupIntent            interface{}         `json:"pending_setup_intent"`
	PendingUpdate                 interface{}         `json:"pending_update"`
	Plan                          SubPlan             `json:"plan"`
	Quantity                      int                 `json:"quantity"`
	Schedule                      interface{}         `json:"schedule"`
	StartDate                     int64               `json:"start_date"`
	Status                        string              `json:"status"`
	TestClock                     interface{}         `json:"test_clock"`
	TransferData                  interface{}         `json:"transfer_data"`
	TrialEnd                      interface{}         `json:"trial_end"`
	TrialSettings                 TrialSettings       `json:"trial_settings"`
	TrialStart                    interface{}         `json:"trial_start"`
}

type SubAutomaticTax struct {
	DisabledReason interface{} `json:"disabled_reason"`
	Enabled        bool        `json:"enabled"`
	Liability      interface{} `json:"liability"`
}

type CancellationDetails struct {
	Comment  interface{} `json:"comment"`
	Feedback interface{} `json:"feedback"`
	Reason   string      `json:"reason"`
}

type InvoiceSettings struct {
	AccountTaxIDs interface{} `json:"account_tax_ids"`
	Issuer        Issuer      `json:"issuer"`
}

type Issuer struct {
	Type string `json:"type"`
}

type Items struct {
	Object     string             `json:"object"`
	Data       []SubscriptionItem `json:"data"`
	HasMore    bool               `json:"has_more"`
	TotalCount int                `json:"total_count"`
	URL        string             `json:"url"`
}

type SubscriptionItem struct {
	ID                string            `json:"id"`
	Object            string            `json:"object"`
	BillingThresholds interface{}       `json:"billing_thresholds"`
	Created           int64             `json:"created"`
	Discounts         []interface{}     `json:"discounts"`
	Metadata          map[string]string `json:"metadata"`
	Plan              SubPlan           `json:"plan"`
	Price             Price             `json:"price"`
	Quantity          int               `json:"quantity"`
	Subscription      string            `json:"subscription"`
	TaxRates          []interface{}     `json:"tax_rates"`
}

type SubPlan struct {
	ID              string                 `json:"id"`
	Object          string                 `json:"object"`
	Active          bool                   `json:"active"`
	AggregateUsage  interface{}            `json:"aggregate_usage"`
	Amount          int                    `json:"amount"`
	AmountDecimal   string                 `json:"amount_decimal"`
	BillingScheme   string                 `json:"billing_scheme"`
	Created         int64                  `json:"created"`
	Currency        string                 `json:"currency"`
	Interval        string                 `json:"interval"`
	IntervalCount   int                    `json:"interval_count"`
	Livemode        bool                   `json:"livemode"`
	Metadata        map[string]interface{} `json:"metadata"`
	Meter           interface{}            `json:"meter"`
	Nickname        interface{}            `json:"nickname"`
	Product         string                 `json:"product"`
	TiersMode       interface{}            `json:"tiers_mode"`
	TransformUsage  interface{}            `json:"transform_usage"`
	TrialPeriodDays interface{}            `json:"trial_period_days"`
	UsageType       string                 `json:"usage_type"`
}

type Price struct {
	ID                string                 `json:"id"`
	Object            string                 `json:"object"`
	Active            bool                   `json:"active"`
	BillingScheme     string                 `json:"billing_scheme"`
	Created           int64                  `json:"created"`
	Currency          string                 `json:"currency"`
	CustomUnitAmount  interface{}            `json:"custom_unit_amount"`
	Livemode          bool                   `json:"livemode"`
	LookupKey         interface{}            `json:"lookup_key"`
	Metadata          map[string]interface{} `json:"metadata"`
	Nickname          interface{}            `json:"nickname"`
	Product           string                 `json:"product"`
	Recurring         Recurring              `json:"recurring"`
	TaxBehavior       string                 `json:"tax_behavior"`
	TiersMode         interface{}            `json:"tiers_mode"`
	TransformQuantity interface{}            `json:"transform_quantity"`
	Type              string                 `json:"type"`
	UnitAmount        int                    `json:"unit_amount"`
	UnitAmountDecimal string                 `json:"unit_amount_decimal"`
}

type Recurring struct {
	AggregateUsage  interface{} `json:"aggregate_usage"`
	Interval        string      `json:"interval"`
	IntervalCount   int         `json:"interval_count"`
	Meter           interface{} `json:"meter"`
	TrialPeriodDays interface{} `json:"trial_period_days"`
	UsageType       string      `json:"usage_type"`
}

type PaymentSettings struct {
	PaymentMethodOptions     SubPaymentMethodOptions `json:"payment_method_options"`
	PaymentMethodTypes       interface{}             `json:"payment_method_types"`
	SaveDefaultPaymentMethod string                  `json:"save_default_payment_method"`
}

type SubPaymentMethodOptions struct {
	ACSSDebit       interface{} `json:"acss_debit"`
	Bancontact      interface{} `json:"bancontact"`
	Card            Card        `json:"card"`
	CustomerBalance interface{} `json:"customer_balance"`
	Konbini         interface{} `json:"konbini"`
	SepaDebit       interface{} `json:"sepa_debit"`
	UsBankAccount   interface{} `json:"us_bank_account"`
}

type Card struct {
	Network             interface{} `json:"network"`
	RequestThreeDSecure string      `json:"request_three_d_secure"`
}

type TrialSettings struct {
	EndBehavior EndBehavior `json:"end_behavior"`
}

type EndBehavior struct {
	MissingPaymentMethod string `json:"missing_payment_method"`
}

type SubRequest struct {
	ID             string      `json:"id"`
	IdempotencyKey interface{} `json:"idempotency_key"`
}
