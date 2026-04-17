package main

import (
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	svc "github.com/romart333/my-go-microservices/payment/pkg/service"
	paymentv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/payment/v1"
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
	paymentv1.RegisterPaymentServiceServer(grpcServer, &svc.PaymentServer{})

	// Включаем reflection для postman/grpcurl
	reflection.Register(grpcServer)

	go func() {
		slog.Info("запуск PaymentService", "адрес", grpcAddress)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("ошибка запуска сервера", "error", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	slog.Info("получен сигнал остановки, завершаем работу")
	grpcServer.GracefulStop()
	slog.Info("grpc server остановлен")
}
