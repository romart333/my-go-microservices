package converter

import (
	"fmt"

	"github.com/google/uuid"

	errs "github.com/romart333/my-go-microservices/order/internal/errors"
	"github.com/romart333/my-go-microservices/order/internal/model"
	"github.com/romart333/my-go-microservices/order/internal/service/input"
	orderv1 "github.com/romart333/my-go-microservices/shared/pkg/openapi/order/v1"
)

func CreateOrderRequestToInput(req *orderv1.CreateOrderRequest) input.CreateOrderInput {
	result := input.CreateOrderInput{
		HullUUID:   req.GetHullUUID(),
		EngineUUID: req.GetEngineUUID(),
	}
	result.ShieldUUID = dtoToUuid(req.GetShieldUUID())
	result.WeaponUUID = dtoToUuid(req.GetWeaponUUID())
	return result
}

func dtoToUuid(dtoUuid orderv1.OptNilUUID) *uuid.UUID {
	if !dtoUuid.IsNull() && dtoUuid.IsSet() {
		result := dtoUuid.Value
		return &result
	}
	return nil
}

func PaymentMethodToModel(pm orderv1.PaymentMethod) model.PaymentMethod {
	switch pm {
	case orderv1.PaymentMethodCARD:
		return model.PaymentMethodCard
	case orderv1.PaymentMethodSBP:
		return model.PaymentMethodSbp
	case orderv1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCreditCard
	case orderv1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnspecified
	}
}

func OrderToDto(order model.Order) (*orderv1.OrderDto, error) {
	var paymentMethod orderv1.OptNilPaymentMethod
	if order.PaymentMethod != nil {
		pm, err := paymentMethodToDto(*order.PaymentMethod)
		if err != nil {
			return &orderv1.OrderDto{}, fmt.Errorf("конвертировать метод оплаты %s в dto: %w", *order.PaymentMethod, err)
		}
		paymentMethod = orderv1.NewOptNilPaymentMethod(pm)
	}

	status, err := statusToDto(order.Status)
	if err != nil {
		return &orderv1.OrderDto{}, fmt.Errorf("конвертировать статус %s в dto: %w", order.Status, err)
	}

	var transactionUUID orderv1.OptNilUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderv1.NewOptNilUUID(*order.TransactionUUID)
	}

	var hullUUID uuid.UUID
	var engineUUID uuid.UUID
	var shieldUUID *uuid.UUID
	var weaponUUID *uuid.UUID
	for _, item := range order.Items {
		switch item.PartType {
		case model.PartTypeHull:
			hullUUID = item.PartUUID
		case model.PartTypeEngine:
			engineUUID = item.PartUUID
		case model.PartTypeShield:
			shieldUUID = &item.PartUUID
		case model.PartTypeWeapon:
			weaponUUID = &item.PartUUID
		default:
			return &orderv1.OrderDto{}, fmt.Errorf("конвертировать тип детали %s в dto: %w", item.PartType, errs.ErrPartUnknownType)
		}
	}

	dto := &orderv1.OrderDto{
		OrderUUID:       (order.UUID),
		TotalPrice:      order.TotalPrice(),
		Status:          status,
		HullUUID:        hullUUID,
		EngineUUID:      engineUUID,
		ShieldUUID:      modelToUUID(shieldUUID),
		WeaponUUID:      modelToUUID(weaponUUID),
		CreatedAt:       order.CreatedAt,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
	}
	return dto, nil
}

func modelToUUID(u *uuid.UUID) orderv1.OptNilUUID {
	if u != nil {
		return orderv1.NewOptNilUUID(*u)
	}
	return orderv1.OptNilUUID{}
}

func paymentMethodToDto(paymentMethod model.PaymentMethod) (orderv1.PaymentMethod, error) {
	switch paymentMethod {
	case model.PaymentMethodCard:
		return orderv1.PaymentMethodCARD, nil
	case model.PaymentMethodSbp:
		return orderv1.PaymentMethodSBP, nil
	case model.PaymentMethodCreditCard:
		return orderv1.PaymentMethodCREDITCARD, nil
	case model.PaymentMethodInvestorMoney:
		return orderv1.PaymentMethodINVESTORMONEY, nil
	default:
		return "", errs.ErrPaymentInvalidMethod
	}
}

func statusToDto(status model.OrderStatus) (orderv1.OrderStatus, error) {
	switch status {
	case model.OrderStatusPendingPayment:
		return orderv1.OrderStatusPENDINGPAYMENT, nil
	case model.OrderStatusPAID:
		return orderv1.OrderStatusPAID, nil
	case model.OrderStatusCancelled:
		return orderv1.OrderStatusCANCELLED, nil
	}
	return "", errs.ErrOrderUnknownStatus
}
