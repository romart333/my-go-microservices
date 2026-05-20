package v1

import (
	"context"
	"fmt"
	"log/slog"

	converter "github.com/romart333/my-go-microservices/payment/internal/api/payment/converter"
	paymentv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/payment/v1"
)

// PayOrder обрабатывает оплату заказа
func (s *server) PayOrder(
	ctx context.Context,
	req *paymentv1.PayOrderRequest,
) (*paymentv1.PayOrderResponse, error) {
	in, err := converter.PayRequestToInput(req)
	if err != nil {
		slog.ErrorContext(ctx, "оплатить заказ: разбор запроса", "order_uuid", req.GetOrderUuid(), "error", err)
		return nil, fmt.Errorf("разобрать запрос оплаты: %w", err)
	}

	transactionUUID, err := s.paymentService.Pay(ctx, in)
	if err != nil {
		slog.ErrorContext(ctx, "оплатить заказ", "order_uuid", in.OrderUUID, "error", err)
		return nil, err
	}

	return &paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID.String(),
	}, nil
}
