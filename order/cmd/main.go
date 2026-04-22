package main

import (
	"context"
	"errors"
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

	orderHandler "github.com/romart333/my-go-microservices/order/pkg/handler"
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
	middlewareTimeout = 10 * time.Second
)

func main() {
	inventoryConn, err := grpc.NewClient(inventoryServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    grpcKeepaliveTime,
			Timeout: grpcKeepaliveTimeout,
		}),
	)
	if err != nil {
		slog.Error("не удалось подключиться к InventoryService", "error", err)
		return
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
		slog.Error("не удалось подключиться к PaymentService", "error", err)
		return
	}
	defer func() {
		if err := paymentConn.Close(); err != nil {
			slog.Error("ошибка закрытия соединения с payment service", "error", err)
		}
	}()

	// Создаём хранилище и обработчик
	store := orderHandler.NewOrderStore()
	h := orderHandler.NewOrderHandler(
		inventoryv1.NewInventoryServiceClient(inventoryConn),
		paymentv1.NewPaymentServiceClient(paymentConn),
		store,
	)

	// Создать OpenAPI сервер
	router, err := orderHandler.SetupServer(h)
	if err != nil {
		slog.Error("ошибка создания сервера OpenAPI", "error", err)
		return
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
		slog.Error("ошибка завершения работы OrderService", "error", err)
	}
	slog.Info("OrderService завершен")
}
