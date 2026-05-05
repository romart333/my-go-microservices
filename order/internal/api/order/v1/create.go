package orderv1

import (
	"context"

	"github.com/google/uuid"

	"github.com/romart333/my-go-microservices/order/internal/converter"
	orderv1 "github.com/romart333/my-go-microservices/shared/pkg/openapi/order/v1"
)

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	modelReq := converter.CreateOrderRequestToModel(req)
	order, err := h.orderService.Create(ctx, modelReq)
	if err != nil {
		return handleCreateOrderError(err)
	}
	return &orderv1.CreateOrderResponse{
		OrderUUID:  uuid.MustParse(order.UUID),
		TotalPrice: order.TotalPrice,
	}, nil
}
