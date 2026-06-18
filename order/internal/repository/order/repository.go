package order

import (
	"sync"

	"github.com/google/uuid"

	"github.com/romart333/my-go-microservices/order/internal/repository/record"
)

// repository — хранилище заказов (in-memory).
type repository struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]record.Order
}

// NewOrderRepository создаёт новое пустое хранилище заказов.
func NewOrderRepository() *repository {
	return &repository{
		orders: make(map[uuid.UUID]record.Order),
	}
}
