package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	errs "github.com/romart333/my-go-microservices/payment/internal/errors"
	"github.com/romart333/my-go-microservices/payment/internal/model"
)

func (s *PaymentService) Pay(ctx context.Context, req model.PayRequest) (string, error) {
	if req.OrderUUID == "" {
		return "", errs.ErrInvalidOrderUUID
	}

	if _, err := uuid.Parse(req.OrderUUID); err != nil {
		return "", errs.ErrInvalidOrderUUID
	}

	if !req.PaymentMethod.IsValid() {
		return "", errs.ErrInvalidPaymentMethod
	}

	transactionUUID := uuid.New().String()

	slog.InfoContext(ctx, "оплата прошла успешно", "order_uuid", req.OrderUUID, "payment_method", req.PaymentMethod)

	return transactionUUID, nil
}
