package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/romart333/my-go-microservices/order/pkg/app"
	inventoryv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/payment/v1"
)

const (
	inventoryServiceAddress = "localhost:50051"
	paymentServiceAddress   = "localhost:50052"
)

const (
	grpcKeepaliveTime    = 10 * time.Second
	grpcKeepaliveTimeout = 3 * time.Second
)

const (
	httpPort = "8080"

	readHeaderTimeout = 5 * time.Second
	readTimeout       = 15 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	if err := run(); err != nil {
		slog.Error("OrderService завершился с ошибкой", "error", err)
		os.Exit(1)
	}
}

func run() error {
	inventoryConn, err := grpc.NewClient(inventoryServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    grpcKeepaliveTime,
			Timeout: grpcKeepaliveTimeout,
		}),
	)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к InventoryService: %w", err)
	}
	defer func() {
		if err := inventoryConn.Close(); err != nil {
			slog.Error("ошибка закрытия соединения с inventory service", "error", err)
		}
	}()

	paymentConn, err := grpc.NewClient(paymentServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    grpcKeepaliveTime,
			Timeout: grpcKeepaliveTimeout,
		}),
	)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к PaymentService: %w", err)
	}
	defer func() {
		if err := paymentConn.Close(); err != nil {
			slog.Error("ошибка закрытия соединения с payment service", "error", err)
		}
	}()

	inventoryClient := inventoryv1.NewInventoryServiceClient(inventoryConn)
	paymentClient := paymentv1.NewPaymentServiceClient(paymentConn)

	router, err := app.NewHTTPHandler(&inventoryClient, &paymentClient)
	if err != nil {
		return fmt.Errorf("ошибка создания сервера OpenAPI: %w", err)
	}

	orderServer := &http.Server{
		Addr:              net.JoinHostPort(":", httpPort),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	go func() {
		slog.Info("запуск OrderService", "port", httpPort)
		if err := orderServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("ошибка запуска сервера", "error", err)
			cancel()
		}
	}()

	<-ctx.Done()
	slog.Info("завершаем работу OrderService ")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	if err := orderServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("ошибка завершения работы OrderService: %w", err)
	}
	slog.Info("OrderService завершен")
	return nil
}
