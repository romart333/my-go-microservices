package orderv1

import (
	"context"
	"log/slog"

	orderv1 "github.com/romart333/my-go-microservices/shared/pkg/openapi/order/v1"
)

func (h *OrderHandler) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	if err := h.orderService.Cancel(ctx, params.OrderUUID.String()); err != nil {
		slog.ErrorContext(ctx, "отмена заказа", "error", err)
		return handleCancelOrderError(err)
	}
	return &orderv1.CancelOrderResponse{}, nil
}
