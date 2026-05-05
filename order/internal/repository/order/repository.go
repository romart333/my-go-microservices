package order

import (
	"sync"

	"github.com/romart333/my-go-microservices/order/internal/repository/record"
)

// OrderStore — хранилище заказов (in-memory).
type OrderStore struct {
	mu     sync.RWMutex
	orders map[string]record.Order
}

// NewOrderStore создаёт новое пустое хранилище заказов.
func NewOrderStore() *OrderStore {
	return &OrderStore{
		orders: make(map[string]record.Order),
	}
}
