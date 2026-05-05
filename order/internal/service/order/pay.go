package order

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/romart333/my-go-microservices/order/internal/errors"
	"github.com/romart333/my-go-microservices/order/internal/model"
)

func (s *OrderService) Pay(ctx context.Context, uuid string, method model.PaymentMethod) (string, error) {
	if uuid == "" {
		return "", errs.ErrOrderIdRequired
	}

	order, err := s.orderRepository.Get(ctx, uuid)
	if err != nil {
		return "", fmt.Errorf("получить заказ: %w", err)
	}

	if order.Status == model.OrderStatusCANCELLED {
		return "", errs.ErrOrderCancelled
	}
	if order.Status == model.OrderStatusPAID {
		return "", errs.ErrOrderAlreadyPaid
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, order.UUID, method)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return "", errs.ErrInvalidParameters
			case codes.NotFound:
				return "", errs.ErrPaymentNotFound
			}
		}
		return "", fmt.Errorf("оплатить заказ: %w", err)
	}

	order.Status = model.OrderStatusPAID
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = &method
	updateErr := s.orderRepository.Update(ctx, order)
	if updateErr != nil {
		return "", fmt.Errorf("обновить заказ: %w", updateErr)
	}

	return transactionUUID, nil
}
