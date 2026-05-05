package order

import (
	"context"
	"fmt"

	errs "github.com/romart333/my-go-microservices/order/internal/errors"
	"github.com/romart333/my-go-microservices/order/internal/model"
)

func (s *OrderService) Cancel(ctx context.Context, uuid string) error {
	if uuid == "" {
		return errs.ErrOrderIdRequired
	}

	order, err := s.orderRepository.Get(ctx, uuid)
	if err != nil {
		return fmt.Errorf("получить заказ: %w", err)
	}

	switch order.Status {
	case model.OrderStatusPENDINGPAYMENT:
		order.Status = model.OrderStatusCANCELLED
		err = s.orderRepository.Update(ctx, order)
		if err != nil {
			return fmt.Errorf("обновить заказ: %w", err)
		}
		return nil
	case model.OrderStatusPAID:
		return errs.ErrOrderAlreadyPaid
	case model.OrderStatusCANCELLED:
		return errs.ErrOrderCancelled
	}

	return errs.ErrOrderUnknownStatus
}
