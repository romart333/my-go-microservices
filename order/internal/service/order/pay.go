package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/romart333/my-go-microservices/order/internal/errors"
	"github.com/romart333/my-go-microservices/order/internal/model"
)

func (s *service) Pay(ctx context.Context, id uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	order, err := s.orderRepository.Get(ctx, id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("получить заказ: %w", err)
	}

	if order.Status == model.OrderStatusCancelled {
		return uuid.Nil, errs.ErrOrderCancelled
	}
	if order.Status == model.OrderStatusPAID {
		return uuid.Nil, errs.ErrOrderAlreadyPaid
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, order.UUID, method)
	if err != nil {
		return uuid.Nil, fmt.Errorf("оплатить заказ: %w", err)
	}

	order.Status = model.OrderStatusPAID
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = &method
	updateErr := s.orderRepository.Update(ctx, order)
	if updateErr != nil {
		return uuid.Nil, fmt.Errorf("обновить заказ: %w", updateErr)
	}

	return transactionUUID, nil
}
