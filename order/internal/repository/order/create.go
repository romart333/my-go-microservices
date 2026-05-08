package order

import (
	"context"

	"github.com/romart333/my-go-microservices/order/internal/model"
	"github.com/romart333/my-go-microservices/order/internal/repository/converter"
)

func (s *OrderStore) Create(_ context.Context, order model.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[order.UUID] = converter.OrderToRecord(order)

	return nil
}
