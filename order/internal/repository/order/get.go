package order

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/romart333/my-go-microservices/order/internal/errors"
	"github.com/romart333/my-go-microservices/order/internal/model"
	"github.com/romart333/my-go-microservices/order/internal/repository/converter"
)

func (s *repository) Get(_ context.Context, uuid uuid.UUID) (model.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[uuid]
	if !ok {
		return model.Order{}, errs.ErrOrderNotFound
	}

	return converter.OrderToModel(order), nil
}
