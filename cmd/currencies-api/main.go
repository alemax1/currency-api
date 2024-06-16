package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alemax1/currencies-api/config"
	forex "github.com/alemax1/currencies-api/internal/currency/adapter/forexApi"
	"github.com/alemax1/currencies-api/internal/currency/adapter/postgres"
	"github.com/alemax1/currencies-api/internal/currency/delivery/http/handler"
	"github.com/alemax1/currencies-api/internal/currency/service"
	"github.com/alemax1/currencies-api/pkg/logger"
	"github.com/alemax1/currencies-api/pkg/pgdb"

	"github.com/gofiber/fiber/v3"
	"github.com/spf13/cobra"
)

const configFlagName = "config"

//	@title			Currencies API
//	@version		1.0
//	@description	API for get currencies rate.

// @host		localhost:3000
// @BasePath	/api/v1
func main() {
	rootCmd := &cobra.Command{
		Use:   "api",
		Short: "currencies api",
		Run: func(cmd *cobra.Command, args []string) {
			cfgPath, err := cmd.Flags().GetString(configFlagName)
			if err != nil {
				log.Fatalf("get flag value: %v", err)
			}

			run(cfgPath)
		},
	}

	rootCmd.Flags().StringP(configFlagName, "c", ".env", "config file path")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func run(cfgPath string) {
	l, err := logger.New()
	if err != nil {
		log.Fatalf("init logger: %v", err)
	}

	cfg, err := config.New(cfgPath)
	if err != nil {
		l.Fatal().Msgf("init config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Postgres.PingTimeout)
	defer cancel()
	db, err := pgdb.Open(ctx, cfg.Postgres.ToDSN())
	if err != nil {
		l.Fatal().Msgf("open db conn: %v", err)
	}

	executor := postgres.NewExecutor(db)
	currencyRepo := postgres.NewCurrency(executor)
	forexApi := forex.New(cfg.CurrenciesAPI)

	service := service.New(currencyRepo, forexApi, l)

	app := fiber.New()
	app.Use(handler.TimeoutMiddleware(cfg.Handler.RequestTimeout))

	handler := handler.New(service, l)
	handler.InitRoutes(app)

	go func() {
		if err := app.Listen(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
			l.Fatal().Msgf("start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-quit
	l.Info().Msg("shutting down server")

	if err := app.Shutdown(); err != nil {
		l.Fatal().Msgf("shutdown server: %v", err)
	}

	l.Info().Msg("server exiting")
}
