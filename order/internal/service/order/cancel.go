package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/romart333/my-go-microservices/order/internal/errors"
	"github.com/romart333/my-go-microservices/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, uuid uuid.UUID) error {
	order, err := s.orderRepository.Get(ctx, uuid)
	if err != nil {
		return fmt.Errorf("получить заказ: %w", err)
	}

	switch order.Status {
	case model.OrderStatusPendingPayment:
		order.Status = model.OrderStatusCancelled
		err = s.orderRepository.Update(ctx, order)
		if err != nil {
			return fmt.Errorf("обновить заказ: %w", err)
		}
		return nil
	case model.OrderStatusPAID:
		return errs.ErrOrderAlreadyPaid
	case model.OrderStatusCancelled:
		return errs.ErrOrderCancelled
	}

	return errs.ErrOrderUnknownStatus
}
