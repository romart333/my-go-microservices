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

	"github.com/romart333/my-go-microservices/payment/pkg/app"
)

const grpcAddress = "localhost:50052"

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
		append(
			app.Interceptors(),
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
		)...,
	)
	app.RegisterServices(grpcServer)

	// Включаем reflection для postman/grpcurl
	reflection.Register(grpcServer)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	go func() {
		slog.Info("запуск PaymentService", "адрес", grpcAddress)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("ошибка запуска сервера", "error", err)
			cancel()
		}
	}()

	<-ctx.Done()
	slog.Info("завершаем работу PaymentService")
	grpcServer.GracefulStop()
	slog.Info("PaymentService остановлен")
}
