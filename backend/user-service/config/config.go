package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServiceName   string
	ServiceID     string
	Host          string
	GRPCPort      int
	MySQLDSN      string
	ConsulAddress string
	JWTSecret     string
	TokenTTL      time.Duration
}

func Load() Config {
	port := envInt("USER_SERVICE_PORT", 9001)
	return Config{
		ServiceName:   "user-service",
		ServiceID:     env("USER_SERVICE_ID", "user-service-1"),
		Host:          env("SERVICE_HOST", "127.0.0.1"),
		GRPCPort:      port,
		MySQLDSN:      env("MYSQL_DSN", "root:123456@tcp(127.0.0.1:3306)/micro_admin?charset=utf8mb4&parseTime=true&loc=Local"),
		ConsulAddress: env("CONSUL_ADDR", "127.0.0.1:8500"),
		JWTSecret:     env("JWT_SECRET", "micro-admin-secret"),
		TokenTTL:      time.Duration(envInt("JWT_TTL_HOURS", 24)) * time.Hour,
	}
}

func env(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func envInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
