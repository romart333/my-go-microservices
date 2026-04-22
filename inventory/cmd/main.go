package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	svc "github.com/romart333/my-go-microservices/inventory/pkg/app/service"
	inventoryv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/inventory/v1"
)

const grpcAddress = "localhost:50051"

const (
	grpcMaxConnectionIdle     = 15 * time.Minute
	grpcMaxConnectionAge      = 30 * time.Minute
	grpcMaxConnectionAgeGrace = 5 * time.Second
	grpcKeepaliveTime         = 5 * time.Minute
	grpcKeepaliveTimeout      = 1 * time.Second
	grpcMinPingInterval       = 5 * time.Minute
)

func main() {
	//nolint:noctx // Контекст здесь не нужен: GracefulStop() сам закроет listener и прервёт Accept()
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		slog.Error("не удалось создать listener", "error", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:                  grpcKeepaliveTime,
			Timeout:               grpcKeepaliveTimeout,
			MaxConnectionIdle:     grpcMaxConnectionIdle,
			MaxConnectionAge:      grpcMaxConnectionAge,
			MaxConnectionAgeGrace: grpcMaxConnectionAgeGrace,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcMinPingInterval,
			PermitWithoutStream: true,
		}),
	)

	inventoryv1.RegisterInventoryServiceServer(grpcServer, svc.NewInventoryServer())

	// Включаем reflection для postman/grpcurl
	reflection.Register(grpcServer)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		slog.Info("запуск InventoryService", "адрес", grpcAddress)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("ошибка запуска InventoryService", "error", err)
			cancel()
		}
	}()

	<-ctx.Done()

	slog.Info("остановка InventoryService")

	grpcServer.GracefulStop()
	slog.Info("InventoryService остановлен")
}
