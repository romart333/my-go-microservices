package app

import (
	"fmt"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	api "github.com/romart333/my-go-microservices/order/internal/api/order/v1"
	inventoryapi "github.com/romart333/my-go-microservices/order/internal/client/grpc/inventory/v1"
	paymentapi "github.com/romart333/my-go-microservices/order/internal/client/grpc/payment/v1"
	orderRepo "github.com/romart333/my-go-microservices/order/internal/repository/order"
	orderService "github.com/romart333/my-go-microservices/order/internal/service/order"
	orderv1 "github.com/romart333/my-go-microservices/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/payment/v1"
)

const (
	middlewareTimeout = 10 * time.Second
)

func registerServices(inventoryClient *inventoryv1.InventoryServiceClient, paymentClient *paymentv1.PaymentServiceClient) *api.OrderHandler {
	inventoryApiClient := inventoryapi.NewInventoryClient(*inventoryClient)
	paymentApiClient := paymentapi.NewPaymentClient(*paymentClient)

	repo := orderRepo.NewOrderRepository()
	service := orderService.NewOrderService(inventoryApiClient, paymentApiClient, repo)
	return api.NewOrderHandler(service)
}

func NewHTTPHandler(inventoryClient *inventoryv1.InventoryServiceClient, paymentClient *paymentv1.PaymentServiceClient) (chi.Router, error) {
	h := registerServices(inventoryClient, paymentClient)
	server, err := orderv1.NewServer(h, orderv1.WithErrorHandler(api.ErrorHandler))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания сервера OpenAPI: %w", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(middlewareTimeout))
	r.Use(middleware.Compress(5))
	r.Handle("/api/*", server)
	return r, nil
}
