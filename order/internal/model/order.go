package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	OrderStatusUNSPECIFIED    OrderStatus = "UNSPECIFIED"
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPAID           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

const (
	PaymentMethodUnspecified   PaymentMethod = "UNSPECIFIED"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSbp           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

type OrderStatus string

type PaymentMethod string

// Order представляет заказ на постройку космического корабля.
type Order struct {
	UUID            uuid.UUID
	Items           []OrderItem
	TransactionUUID *uuid.UUID
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
}

// TotalPrice возвращает сумму цен всех позиций заказа.
func (o Order) TotalPrice() int64 {
	var total int64
	for _, item := range o.Items {
		total += item.Price
	}
	return total
}

type OrderItem struct {
	PartUUID uuid.UUID
	PartType PartType
	Price    int64
}

func (o OrderStatus) IsValid() bool {
	switch o {
	case OrderStatusPendingPayment, OrderStatusPAID, OrderStatusCancelled:
		return true
	}
	return false
}

func (o PaymentMethod) IsValid() bool {
	switch o {
	case PaymentMethodCard, PaymentMethodSbp, PaymentMethodCreditCard, PaymentMethodInvestorMoney:
		return true
	}

	return false
}
