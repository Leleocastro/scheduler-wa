package main

import (
	"complete-api/internal/adapters/handlers/checkouthdl"
	"complete-api/internal/adapters/handlers/gatewayhdl"
	"complete-api/internal/adapters/handlers/scheduleshdl"
	"complete-api/internal/adapters/handlers/statshdl"
	"complete-api/internal/adapters/handlers/validatorhdl"
	"complete-api/internal/adapters/repositories/api_gateway"
	"complete-api/internal/adapters/repositories/payment"
	"complete-api/internal/adapters/repositories/redis"
	"complete-api/internal/adapters/repositories/stats"
	"complete-api/internal/core/services/checkoutsrv"
	"complete-api/internal/core/services/gatewaysrv"
	"complete-api/internal/core/services/paymentsrv"
	"complete-api/internal/core/services/schedulessrv"
	"complete-api/internal/core/services/statssrv"
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Erro ao carregar o arquivo .env:", err)
	}
	// Inicializa a conexão com o banco de dados
	initDB()

	kongRepo := api_gateway.New(os.Getenv("KONG_ADMIN_URL"))
	stripeRepo := payment.New(os.Getenv("STRIPE_SECRET_KEY"))
	prometheusRepo := stats.New(os.Getenv("PROMETHEUS_URL"))

	paymentSrv := paymentsrv.New(stripeRepo)

	validatorHandler := validatorhdl.NewValidatorHandler(paymentSrv)

	// Cria o router Gin
	router := gin.Default()

	router.RedirectTrailingSlash = true

	checkout := router.Group("/checkout")
	{
		checkoutSrv := checkoutsrv.New(kongRepo)

		checkoutHandler := checkouthdl.NewHTTPHandler(checkoutSrv, paymentSrv)

		checkout.POST("/webhook", validatorHandler.ValidateSignature, checkoutHandler.Checkout)
		checkout.POST("/cancel", validatorHandler.ValidateSignature, checkoutHandler.CancelSubscription)
	}

	kong := router.Group("/gateway")
	{
		gatewaySrv := gatewaysrv.New(kongRepo)

		gatewayHandler := gatewayhdl.NewHTTPHandler(gatewaySrv)

		kong.POST("/consumer", gatewayHandler.CreateConsumer)

		kong.GET("/api-key", gatewayHandler.GetAPIKey)
	}

	api := router.Group("/api")
	{
		statsSrv := statssrv.New(prometheusRepo)

		statsHandler := statshdl.NewHTTPHandler(statsSrv)

		api.GET("/usage/:consumer", statsHandler.GetUsageByConsumer)
	}

	schedules := router.Group("/schedules")
	{
		redisRepo := redis.New(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PASSWORD"))

		scheduleSrv := schedulessrv.New(redisRepo)

		go scheduleSrv.Run()

		scheduleHandler := scheduleshdl.NewHTTPHandler(scheduleSrv)

		schedules.POST("/", scheduleHandler.CreateSchedule)
		schedules.GET("/", scheduleHandler.GetSchedules)
		schedules.PUT("/:scheduleID", scheduleHandler.UpdateSchedule)
		schedules.DELETE("/:scheduleID", scheduleHandler.DeleteSchedule)
	}

	router.Run(":6000")
}
