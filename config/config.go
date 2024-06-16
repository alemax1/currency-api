package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Postgres         Postgres
	CurrenciesAPI    CurrenciesAPI
	Handler          Handler
	CurrenciesWorker CurrenciesWorker
	Server           Server
}

func New(cfgPath string) (*Config, error) {
	if err := godotenv.Load(cfgPath); err != nil {
		return nil, err
	}

	return &Config{
		Postgres:         newPostgres(),
		CurrenciesAPI:    newCurrenciesAPI(),
		Handler:          newHandler(),
		CurrenciesWorker: newCurrenciesWorker(),
		Server:           newServer(),
	}, nil
}

func getDefaultEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func getDefaultIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	val, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal(err)
	}

	return val
}

func getDefaultDurationEnv(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	switch value[len(value)-1] {
	case 'm':
		duration, err := strconv.Atoi(value[:len(value)-1])
		if err != nil {
			log.Fatal(err)
		}

		return time.Duration(duration) * time.Minute
	case 's':
		duration, err := strconv.Atoi(value[:len(value)-1])
		if err != nil {
			log.Fatal(err)
		}

		return time.Duration(duration) * time.Second
	case 'h':
		duration, err := strconv.Atoi(value[:len(value)-1])
		if err != nil {
			log.Fatal(err)
		}

		return time.Duration(duration) * time.Hour
	case 'l':
		duration, err := strconv.Atoi(value[:len(value)-1])
		if err != nil {
			log.Fatal(err)
		}

		return time.Duration(duration) * time.Millisecond
	}

	return defaultValue
}
