package ports

type APIGatewayRepository interface {
	CreateConsumer(username, customID string) error
	RateLimitConsumer(username, route string, rateLimit int) error
	CreateJWTFirebaseConsumer(username string) error
}
