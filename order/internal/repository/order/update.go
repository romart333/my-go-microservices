package order

import (
	"context"

	errs "github.com/romart333/my-go-microservices/order/internal/errors"
	"github.com/romart333/my-go-microservices/order/internal/model"
	"github.com/romart333/my-go-microservices/order/internal/repository/converter"
)

func (s *repository) Update(_ context.Context, order model.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.orders[order.UUID]
	if !ok {
		return errs.ErrOrderNotFound
	}

	s.orders[order.UUID] = converter.OrderToRecord(order)

	return nil
}
