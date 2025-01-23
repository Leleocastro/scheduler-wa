package ports

type APIGatewayRepository interface {
	CreateConsumer(username, customID string) error
	RateLimitConsumer(username, route string, rateLimit int) error
	CreateACL(username, group string) error
	CreateAPIKey(username string) error
}

type PaymentRepository interface {
	ValidateSignature(payload []byte, sigHeader string) error
	GetPlanByPriceID(priceID string) (string, error)
}
