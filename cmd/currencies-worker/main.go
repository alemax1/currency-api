package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alemax1/currencies-api/config"
	worker "github.com/alemax1/currencies-api/internal/currencies-worker"
	forex "github.com/alemax1/currencies-api/internal/currency/adapter/forexApi"
	"github.com/alemax1/currencies-api/internal/currency/adapter/postgres"
	"github.com/alemax1/currencies-api/internal/currency/service"
	"github.com/alemax1/currencies-api/pkg/logger"
	"github.com/alemax1/currencies-api/pkg/pgdb"
	"github.com/spf13/cobra"
)

const configFlagName = "config"

func main() {
	rootCmd := &cobra.Command{
		Use:   "worker",
		Short: "currencies worker",
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
		l.Fatal().Msgf("pg conn open: %v", err)
	}

	executor := postgres.NewExecutor(db)
	currencyRepo := postgres.NewCurrency(executor)

	forexApi := forex.New(cfg.CurrenciesAPI)

	service := service.New(currencyRepo, forexApi, l)

	worker := worker.New(service.Currency, cfg.CurrenciesWorker, l)
	worker.Run()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	<-quit

	worker.Stop()
}
