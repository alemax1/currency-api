package config

import (
	"fmt"
	"time"
)

type Postgres struct {
	Host         string
	Port         int
	Password     string
	Username     string
	DatabaseName string
	PingTimeout  time.Duration
}

func newPostgres() Postgres {
	return Postgres{
		Host:         getDefaultEnv("POSTGRES_HOST", "localhost"),
		Port:         getDefaultIntEnv("POSTGRES_PORT", 5440),
		Password:     getDefaultEnv("POSTGRES_PASSWORD", "123456"),
		Username:     getDefaultEnv("POSTGRES_USERNAME", "postgres"),
		DatabaseName: getDefaultEnv("POSTGRES_DATABASE", "postgres"),
		PingTimeout:  getDefaultDurationEnv("POSTGRES_PING_TIMEOUT", 5*time.Minute),
	}
}

func (p Postgres) ToDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		p.Username, p.Password, p.Host, p.Port, p.DatabaseName)
}
