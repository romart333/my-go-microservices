package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/romart333/my-go-microservices/payment/internal/errors"
	input "github.com/romart333/my-go-microservices/payment/internal/service/input"
)

func (s *service) Pay(ctx context.Context, req input.PayOrderInput) (uuid.UUID, error) {
	if !req.PaymentMethod.IsValid() {
		return uuid.Nil, errs.ErrInvalidPaymentMethod
	}

	transactionUUID := uuid.New()

	slog.InfoContext(ctx, "оплата прошла успешно", "order_uuid", req.OrderUUID, "payment_method", req.PaymentMethod)

	return transactionUUID, nil
}
