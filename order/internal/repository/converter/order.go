package converter

import (
	"github.com/romart333/my-go-microservices/order/internal/model"
	"github.com/romart333/my-go-microservices/order/internal/repository/record"
)

func OrderToRecord(o model.Order) record.Order {
	return record.Order{
		OrderUUID:       o.UUID,
		HullUUID:        o.HullUUID,
		EngineUUID:      o.EngineUUID,
		ShieldUUID:      o.ShieldUUID,
		WeaponUUID:      o.WeaponUUID,
		TotalPrice:      o.TotalPrice,
		TransactionUUID: o.TransactionUUID,
		PaymentMethod:   paymentMethodToRecord(o.PaymentMethod),
		Status:          orderStatusToRecord(o.Status),
		CreatedAt:       o.CreatedAt,
	}
}

func OrderToModel(o record.Order) model.Order {
	return model.Order{
		UUID:            o.OrderUUID,
		PaymentMethod:   paymentMethodToModel(o.PaymentMethod),
		Status:          orderStatusToModel(o.Status),
		CreatedAt:       o.CreatedAt,
		HullUUID:        o.HullUUID,
		EngineUUID:      o.EngineUUID,
		ShieldUUID:      o.ShieldUUID,
		WeaponUUID:      o.WeaponUUID,
		TotalPrice:      o.TotalPrice,
		TransactionUUID: o.TransactionUUID,
	}
}

func paymentMethodToModel(pm *string) *model.PaymentMethod {
	if pm == nil {
		return nil
	}
	switch *pm {
	case "CARD":
		return new(model.PaymentMethodCARD)
	case "SBP":
		return new(model.PaymentMethodSBP)
	case "CREDIT_CARD":
		return new(model.PaymentMethodCREDITCARD)
	case "INVESTOR_MONEY":
		return new(model.PaymentMethodINVESTORMONEY)
	}
	return new(model.PaymentMethodUNSPECIFIED)
}

func paymentMethodToRecord(pm *model.PaymentMethod) *string {
	if pm == nil {
		return nil
	}
	switch *pm {
	case model.PaymentMethodCARD:
		return new("CARD")
	case model.PaymentMethodSBP:
		return new("SBP")
	case model.PaymentMethodCREDITCARD:
		return new("CREDIT_CARD")
	case model.PaymentMethodINVESTORMONEY:
		return new("INVESTOR_MONEY")
	}
	return new("UNSPECIFIED")
}

func orderStatusToModel(os string) model.OrderStatus {
	switch os {
	case "PENDING_PAYMENT":
		return model.OrderStatusPENDINGPAYMENT
	case "PAID":
		return model.OrderStatusPAID
	case "CANCELLED":
		return model.OrderStatusCANCELLED
	}
	return model.OrderStatusUNSPECIFIED
}

func orderStatusToRecord(os model.OrderStatus) string {
	switch os {
	case model.OrderStatusPENDINGPAYMENT:
		return "PENDING_PAYMENT"
	case model.OrderStatusPAID:
		return "PAID"
	case model.OrderStatusCANCELLED:
		return "CANCELLED"
	}
	return "UNSPECIFIED"
}
