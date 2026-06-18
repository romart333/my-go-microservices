package orderv1

import (
	"context"
	"log/slog"

	converter "github.com/romart333/my-go-microservices/order/internal/api/order/converter"
	orderv1 "github.com/romart333/my-go-microservices/shared/pkg/openapi/order/v1"
)

func (h *OrderHandler) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	transactionUUID, err := h.orderService.Pay(ctx, params.OrderUUID, converter.PaymentMethodToModel(req.GetPaymentMethod()))
	if err != nil {
		slog.ErrorContext(ctx, "оплатить заказ", "order_uuid", params.OrderUUID, "error", err)
		return nil, err
	}

	return &orderv1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}
