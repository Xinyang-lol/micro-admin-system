package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPPort      int
	ConsulAddress string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	JWTSecret     string
	RequestTimeout time.Duration
	UploadDir     string
}

func Load() Config {
	return Config{
		HTTPPort:       envInt("GATEWAY_PORT", 8080),
		ConsulAddress:  env("CONSUL_ADDR", "127.0.0.1:8500"),
		RedisAddr:      env("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword:  env("REDIS_PASSWORD", ""),
		RedisDB:        envInt("REDIS_DB", 0),
		JWTSecret:      env("JWT_SECRET", "micro-admin-secret"),
		RequestTimeout: time.Duration(envInt("REQUEST_TIMEOUT_SECONDS", 3)) * time.Second,
		UploadDir:      env("UPLOAD_DIR", "./uploads"),
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
