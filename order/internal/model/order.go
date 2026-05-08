package model

import (
	"time"
)

const (
	OrderStatusUNSPECIFIED    OrderStatus = "UNSPECIFIED"
	OrderStatusPENDINGPAYMENT OrderStatus = "PENDING_PAYMENT"
	OrderStatusPAID           OrderStatus = "PAID"
	OrderStatusCANCELLED      OrderStatus = "CANCELLED"
)

const (
	PaymentMethodUNSPECIFIED   PaymentMethod = "UNSPECIFIED"
	PaymentMethodCARD          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCREDITCARD    PaymentMethod = "CREDIT_CARD"
	PaymentMethodINVESTORMONEY PaymentMethod = "INVESTOR_MONEY"
)

type OrderStatus string

type PaymentMethod string

// Order представляет заказ на постройку космического корабля.
type Order struct {
	UUID            string
	HullUUID        string
	EngineUUID      string
	ShieldUUID      *string
	WeaponUUID      *string
	TotalPrice      int64
	TransactionUUID *string
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
}

func (o OrderStatus) IsValid() bool {
	switch o {
	case OrderStatusPENDINGPAYMENT, OrderStatusPAID, OrderStatusCANCELLED:
		return true
	}
	return false
}

func (o PaymentMethod) IsValid() bool {
	switch o {
	case PaymentMethodCARD, PaymentMethodSBP, PaymentMethodCREDITCARD, PaymentMethodINVESTORMONEY:
		return true
	}

	return false
}
