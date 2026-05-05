package order

import (
	"github.com/romart333/my-go-microservices/order/internal/model"
	"github.com/romart333/my-go-microservices/order/internal/repository/converter"
)

func (s *OrderStore) CreateOrder(order model.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[order.UUID] = converter.OrderToRecord(order)

	return nil
}
