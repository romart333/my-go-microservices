package order

import (
	"context"
	"fmt"

	errs "github.com/romart333/my-go-microservices/order/internal/errors"
	"github.com/romart333/my-go-microservices/order/internal/model"
)

func (s *OrderService) Get(ctx context.Context, uuid string) (model.Order, error) {
	if uuid == "" {
		return model.Order{}, errs.ErrOrderIdRequired
	}

	order, err := s.orderRepository.Get(ctx, uuid)
	if err != nil {
		return model.Order{}, fmt.Errorf("получить заказ: %w", err)
	}

	return order, nil
}
