package converter

import (
	"github.com/romart333/my-go-microservices/order/internal/model"
	"github.com/romart333/my-go-microservices/order/internal/repository/record"
)

func OrderToRecord(o model.Order) record.Order {
	return record.Order{
		OrderUUID:       o.UUID,
		OrderItems:      orderItemsToRecords(o.Items),
		TotalPrice:      o.TotalPrice(),
		TransactionUUID: o.TransactionUUID,
		PaymentMethod:   paymentMethodToRecord(o.PaymentMethod),
		Status:          orderStatusToRecord(o.Status),
		CreatedAt:       o.CreatedAt,
	}
}

func orderItemsToRecords(items []model.OrderItem) []record.OrderItem {
	records := make([]record.OrderItem, 0, len(items))
	for _, item := range items {
		records = append(records, record.OrderItem{
			PartUUID: item.PartUUID,
			PartType: partTypeToRecord(item.PartType),
			Price:    item.Price,
		})
	}
	return records
}

func partTypeToRecord(partType model.PartType) string {
	switch partType {
	case model.PartTypeHull:
		return "HULL"
	case model.PartTypeEngine:
		return "ENGINE"
	case model.PartTypeShield:
		return "SHIELD"
	case model.PartTypeWeapon:
		return "WEAPON"
	}
	return "UNSPECIFIED"
}

func OrderToModel(o record.Order) model.Order {
	return model.Order{
		UUID:            o.OrderUUID,
		PaymentMethod:   paymentMethodToModel(o.PaymentMethod),
		Status:          orderStatusToModel(o.Status),
		CreatedAt:       o.CreatedAt,
		Items:           orderItemsToModels(o.OrderItems),
		TransactionUUID: o.TransactionUUID,
	}
}

func orderItemsToModels(items []record.OrderItem) []model.OrderItem {
	models := make([]model.OrderItem, 0, len(items))
	for _, item := range items {
		models = append(models, model.OrderItem{
			PartUUID: item.PartUUID,
			PartType: partTypeToModel(item.PartType),
			Price:    item.Price,
		})
	}
	return models
}

func partTypeToModel(partType string) model.PartType {
	switch partType {
	case "HULL":
		return model.PartTypeHull
	case "ENGINE":
		return model.PartTypeEngine
	case "SHIELD":
		return model.PartTypeShield
	case "WEAPON":
		return model.PartTypeWeapon
	}
	return model.PartTypeUnspecified
}

func paymentMethodToModel(pm *string) *model.PaymentMethod {
	if pm == nil {
		return nil
	}
	switch *pm {
	case "CARD":
		return new(model.PaymentMethodCard)
	case "SBP":
		return new(model.PaymentMethodSbp)
	case "CREDIT_CARD":
		return new(model.PaymentMethodCreditCard)
	case "INVESTOR_MONEY":
		return new(model.PaymentMethodInvestorMoney)
	}
	return new(model.PaymentMethodUnspecified)
}

func paymentMethodToRecord(pm *model.PaymentMethod) *string {
	if pm == nil {
		return nil
	}
	switch *pm {
	case model.PaymentMethodCard:
		return new("CARD")
	case model.PaymentMethodSbp:
		return new("SBP")
	case model.PaymentMethodCreditCard:
		return new("CREDIT_CARD")
	case model.PaymentMethodInvestorMoney:
		return new("INVESTOR_MONEY")
	}
	return new("UNSPECIFIED")
}

func orderStatusToModel(os string) model.OrderStatus {
	switch os {
	case "PENDING_PAYMENT":
		return model.OrderStatusPendingPayment
	case "PAID":
		return model.OrderStatusPAID
	case "CANCELLED":
		return model.OrderStatusCancelled
	}
	return model.OrderStatusUNSPECIFIED
}

func orderStatusToRecord(os model.OrderStatus) string {
	switch os {
	case model.OrderStatusPendingPayment:
		return "PENDING_PAYMENT"
	case model.OrderStatusPAID:
		return "PAID"
	case model.OrderStatusCancelled:
		return "CANCELLED"
	}
	return "UNSPECIFIED"
}
