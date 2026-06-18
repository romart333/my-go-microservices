package orderv1

import (
	"context"
	"log/slog"

	converter "github.com/romart333/my-go-microservices/order/internal/api/order/converter"
	orderv1 "github.com/romart333/my-go-microservices/shared/pkg/openapi/order/v1"
)

func (h *OrderHandler) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	order, err := h.orderService.Get(ctx, params.OrderUUID)
	if err != nil {
		slog.ErrorContext(ctx, "получить заказ", "order_uuid", params.OrderUUID, "error", err)
		return nil, err
	}
	dto, err := converter.OrderToDto(order)
	if err != nil {
		slog.ErrorContext(ctx, "конвертировать заказ в DTO", "order_uuid", order.UUID, "error", err)
		return nil, err
	}
	return dto, nil
}
