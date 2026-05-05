package orderv1

import (
	"context"

	"github.com/google/uuid"

	"github.com/romart333/my-go-microservices/order/internal/converter"
	orderv1 "github.com/romart333/my-go-microservices/shared/pkg/openapi/order/v1"
)

func (h *OrderHandler) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	pm := req.GetPaymentMethod()
	order, err := h.orderService.Pay(ctx, params.OrderUUID.String(), converter.PaymentMethodToModel(pm))
	if err != nil {
		return handlePayOrderError(err)
	}

	return &orderv1.PayOrderResponse{
		TransactionUUID: uuid.MustParse(order),
	}, nil
}
