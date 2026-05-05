package orderv1

import orderv1 "github.com/romart333/my-go-microservices/shared/pkg/openapi/order/v1"

type OrderHandler struct {
	orderService OrderService
	orderv1.UnimplementedHandler
}

func NewOrderHandler(orderService OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}
