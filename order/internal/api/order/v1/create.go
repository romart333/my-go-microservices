package orderv1

import (
	"context"
	"log/slog"

	converter "github.com/romart333/my-go-microservices/order/internal/api/order/converter"
	orderv1 "github.com/romart333/my-go-microservices/shared/pkg/openapi/order/v1"
)

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	in := converter.CreateOrderRequestToInput(req)
	order, err := h.orderService.Create(ctx, in)
	if err != nil {
		slog.ErrorContext(ctx, "создать заказ", "error", err)
		return nil, err
	}
	return &orderv1.CreateOrderResponse{
		OrderUUID:  order.UUID,
		TotalPrice: order.TotalPrice(),
	}, nil
}
