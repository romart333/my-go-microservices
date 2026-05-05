package order

import (
	inventory "github.com/romart333/my-go-microservices/order/internal/client/grpc/inventory/v1"
	payment "github.com/romart333/my-go-microservices/order/internal/client/grpc/payment/v1"
	orderRepo "github.com/romart333/my-go-microservices/order/internal/repository/order"
)

type OrderService struct {
	inventoryClient *inventory.Client
	paymentClient   *payment.Client
	orderRepository *orderRepo.OrderStore
}

func NewOrderService(inventoryClient *inventory.Client, paymentClient *payment.Client, orderRepository *orderRepo.OrderStore) *OrderService {
	return &OrderService{
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		orderRepository: orderRepository,
	}
}
