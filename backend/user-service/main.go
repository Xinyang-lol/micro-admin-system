package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/signal"
	"syscall"

	"micro-admin-system/backend/common/db"
	"micro-admin-system/backend/common/grpcjson"
	"micro-admin-system/backend/common/registry"
	"micro-admin-system/backend/proto/gen"
	"micro-admin-system/backend/user-service/config"
	"micro-admin-system/backend/user-service/repository"
	usersvc "micro-admin-system/backend/user-service/service"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	database, err := db.OpenMySQL(cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("open mysql failed: %v", err)
	}
	defer database.Close()

	grpcjson.Register()
	server := grpc.NewServer(grpc.ForceServerCodec(grpcjson.Codec{}))
	pb.RegisterUserServiceServer(server, usersvc.New(repository.New(database), cfg.JWTSecret, cfg.TokenTTL))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("listen failed: %v", err)
	}

	consul, err := registry.New(cfg.ConsulAddress)
	if err != nil {
		log.Fatalf("connect consul failed: %v", err)
	}
	if err := consul.Register(ctx, cfg.ServiceName, cfg.ServiceID, cfg.Host, cfg.GRPCPort, []string{"grpc", "user"}); err != nil {
		log.Fatalf("register consul failed: %v", err)
	}
	registry.StartRuntimeReporter(ctx, cfg.ServiceName)

	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()

	log.Printf("%s started at :%d", cfg.ServiceName, cfg.GRPCPort)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("grpc server stopped: %v", err)
	}
}
