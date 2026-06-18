package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/romart333/my-go-microservices/order/internal/model"
)

func (s *service) Get(ctx context.Context, uuid uuid.UUID) (model.Order, error) {
	order, err := s.orderRepository.Get(ctx, uuid)
	if err != nil {
		return model.Order{}, fmt.Errorf("получить заказ: %w", err)
	}

	return order, nil
}
