package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/romart333/my-go-microservices/payment/internal/converter"
	errs "github.com/romart333/my-go-microservices/payment/internal/errors"
	paymentv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/payment/v1"
)

// PayOrder обрабатывает оплату заказа
func (s *PaymentServer) PayOrder(
	ctx context.Context,
	req *paymentv1.PayOrderRequest,
) (*paymentv1.PayOrderResponse, error) {
	model := converter.PayRequestToModel(req)
	transactionUUID, err := s.paymentService.Pay(ctx, model)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidOrderUUID) {
			return nil, status.Error(codes.InvalidArgument, "неверный формат uuid заказа")
		}
		if errors.Is(err, errs.ErrInvalidPaymentMethod) {
			return nil, status.Error(codes.InvalidArgument, "неверный метод оплаты")
		}
		return nil, status.Error(codes.Internal, "ошибка при оплате заказа")
	}

	return &paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
