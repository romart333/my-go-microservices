package orderv1

import (
	"context"
	"net/http"

	"github.com/romart333/my-go-microservices/order/internal/converter"
	orderv1 "github.com/romart333/my-go-microservices/shared/pkg/openapi/order/v1"
)

func (h *OrderHandler) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	order, err := h.orderService.Get(ctx, params.OrderUUID.String())
	if err != nil {
		return handleGetOrderError(err)
	}
	dto, err := converter.OrderToDto(order)
	if err != nil {
		return &orderv1.GetOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "ошибка при получении заказа",
		}, nil
	}
	return dto, nil
}
