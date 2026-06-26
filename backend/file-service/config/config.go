package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServiceName   string
	ServiceID     string
	Host          string
	GRPCPort      int
	MySQLDSN      string
	ConsulAddress string
}

func Load() Config {
	return Config{
		ServiceName:   "file-service",
		ServiceID:     env("FILE_SERVICE_ID", "file-service-1"),
		Host:          env("SERVICE_HOST", "127.0.0.1"),
		GRPCPort:      envInt("FILE_SERVICE_PORT", 9003),
		MySQLDSN:      env("MYSQL_DSN", "root:123456@tcp(127.0.0.1:3306)/micro_admin?charset=utf8mb4&parseTime=true&loc=Local"),
		ConsulAddress: env("CONSUL_ADDR", "127.0.0.1:8500"),
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
