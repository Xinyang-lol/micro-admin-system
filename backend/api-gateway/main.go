package main

import (
	"fmt"
	"log"
	"os"

	"micro-admin-system/backend/api-gateway/client"
	"micro-admin-system/backend/api-gateway/config"
	"micro-admin-system/backend/api-gateway/handler"
	"micro-admin-system/backend/api-gateway/router"
	"micro-admin-system/backend/common/grpcjson"
	redisx "micro-admin-system/backend/common/redis"
	"micro-admin-system/backend/common/registry"
)

func main() {
	cfg := config.Load()
	grpcjson.Register()

	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		log.Fatalf("create upload dir failed: %v", err)
	}

	redisClient, err := redisx.NewClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("connect redis failed: %v", err)
	}
	defer redisClient.Close()

	discovery, err := registry.New(cfg.ConsulAddress)
	if err != nil {
		log.Fatalf("connect consul failed: %v", err)
	}
	clients := client.New(discovery)
	h := handler.New(clients, cfg.RequestTimeout, redisClient, cfg.UploadDir)
	r := router.Setup(h, cfg.JWTSecret, redisClient, clients, cfg.RequestTimeout)

	addr := fmt.Sprintf(":%d", cfg.HTTPPort)
	log.Printf("api-gateway started at %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("gateway stopped: %v", err)
	}
}
