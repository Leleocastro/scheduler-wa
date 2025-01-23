package main

import (
	"complete-api/internal/adapters/handlers/checkouthdl"
	"complete-api/internal/adapters/handlers/gatewayhdl"
	"complete-api/internal/adapters/handlers/validatorhdl"
	"complete-api/internal/adapters/repositories/api_gateway"
	"complete-api/internal/adapters/repositories/payment"
	"complete-api/internal/core/services/checkoutsrv"
	"complete-api/internal/core/services/gatewaysrv"
	"complete-api/internal/core/services/paymentsrv"
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

// Configura a conexão com o banco de dados PostgreSQL
func initDB() {
	var err error
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"))
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Conectado ao banco de dados com sucesso!")
}

func main() {
	// Inicializa a conexão com o banco de dados
	initDB()

	kongRepo := api_gateway.New(os.Getenv("KONG_ADMIN_URL"))
	stripeRepo := payment.New(os.Getenv("STRIPE_SECRET_KEY"), os.Getenv("STRIPE_WEBHOOK_SECRET"))

	paymentSrv := paymentsrv.New(stripeRepo)

	validatorHandler := validatorhdl.NewValidatorHandler(paymentSrv)

	// Cria o router Gin
	router := gin.Default()

	checkout := router.Group("/checkout")
	{
		checkoutSrv := checkoutsrv.New(kongRepo)

		checkoutHandler := checkouthdl.NewHTTPHandler(checkoutSrv, paymentSrv)

		checkout.POST("/webhook", validatorHandler.ValidateSignature, checkoutHandler.Checkout)
	}

	kong := router.Group("/gateway")
	{
		gatewaySrv := gatewaysrv.New(kongRepo)

		gatewayHandler := gatewayhdl.NewHTTPHandler(gatewaySrv)

		kong.POST("/consumer", gatewayHandler.CreateConsumer)

		kong.GET("/api-key", gatewayHandler.GetAPIKey)
	}

	router.Run(":6000")
}
