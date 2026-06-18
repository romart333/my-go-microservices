package record

import (
	"time"

	"github.com/google/uuid"
)

// Order представляет заказ на постройку космического корабля.
type Order struct {
	OrderUUID       uuid.UUID
	OrderItems      []OrderItem
	TotalPrice      int64 // в копейках
	TransactionUUID *uuid.UUID
	PaymentMethod   *string
	Status          string // PENDING_PAYMENT, PAID, CANCELLED
	CreatedAt       time.Time
}

type OrderItem struct {
	PartUUID uuid.UUID
	PartType string
	Price    int64
}
