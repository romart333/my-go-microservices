package converter

import (
	"fmt"

	"github.com/google/uuid"

	errs "github.com/romart333/my-go-microservices/order/internal/errors"
	"github.com/romart333/my-go-microservices/order/internal/model"
	orderv1 "github.com/romart333/my-go-microservices/shared/pkg/openapi/order/v1"
)

func CreateOrderRequestToModel(req *orderv1.CreateOrderRequest) model.CreateOrderRequest {
	result := model.CreateOrderRequest{
		HullUUID:   req.GetHullUUID().String(),
		EngineUUID: req.GetEngineUUID().String(),
	}
	if !req.GetShieldUUID().IsNull() && req.GetShieldUUID().IsSet() {
		shield := req.GetShieldUUID().Value.String()
		result.ShieldUUID = &shield
	}
	if !req.GetWeaponUUID().IsNull() && req.GetWeaponUUID().IsSet() {
		weapon := req.GetWeaponUUID().Value.String()
		result.WeaponUUID = &weapon
	}
	return result
}

func PaymentMethodToModel(pm orderv1.PaymentMethod) model.PaymentMethod {
	switch pm {
	case orderv1.PaymentMethodCARD:
		return model.PaymentMethodCARD
	case orderv1.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case orderv1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCREDITCARD
	case orderv1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodINVESTORMONEY
	default:
		return model.PaymentMethodUNSPECIFIED
	}
}

func OrderToDto(order model.Order) (*orderv1.OrderDto, error) {
	var paymentMethod orderv1.OptNilPaymentMethod
	if order.PaymentMethod != nil {
		pm, err := paymentMethodToDto(*order.PaymentMethod)
		if err != nil {
			return &orderv1.OrderDto{}, fmt.Errorf("convert payment method %s to dto: %w", *order.PaymentMethod, err)
		}
		paymentMethod = orderv1.NewOptNilPaymentMethod(pm)
	}

	status, err := statusToDto(order.Status)
	if err != nil {
		return &orderv1.OrderDto{}, fmt.Errorf("convert status %s to dto: %w", order.Status, err)
	}

	var shieldUUID orderv1.OptNilUUID
	if order.ShieldUUID != nil {
		shieldUUID = orderv1.NewOptNilUUID(uuid.MustParse(*order.ShieldUUID))
	}
	var weaponUUID orderv1.OptNilUUID
	if order.WeaponUUID != nil {
		weaponUUID = orderv1.NewOptNilUUID(uuid.MustParse(*order.WeaponUUID))
	}

	var transactionUUID orderv1.OptNilUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderv1.NewOptNilUUID(uuid.MustParse(*order.TransactionUUID))
	}

	dto := &orderv1.OrderDto{
		OrderUUID:       uuid.MustParse(order.UUID),
		TotalPrice:      order.TotalPrice,
		Status:          status,
		CreatedAt:       order.CreatedAt,
		HullUUID:        uuid.MustParse(order.HullUUID),
		EngineUUID:      uuid.MustParse(order.EngineUUID),
		ShieldUUID:      shieldUUID,
		WeaponUUID:      weaponUUID,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
	}
	return dto, nil
}

func paymentMethodToDto(paymentMethod model.PaymentMethod) (orderv1.PaymentMethod, error) {
	switch paymentMethod {
	case model.PaymentMethodCARD:
		return orderv1.PaymentMethodCARD, nil
	case model.PaymentMethodSBP:
		return orderv1.PaymentMethodSBP, nil
	case model.PaymentMethodCREDITCARD:
		return orderv1.PaymentMethodCREDITCARD, nil
	case model.PaymentMethodINVESTORMONEY:
		return orderv1.PaymentMethodINVESTORMONEY, nil
	default:
		return "", errs.ErrOrderInvalidPaymentMethod
	}
}

func statusToDto(status model.OrderStatus) (orderv1.OrderStatus, error) {
	switch status {
	case model.OrderStatusPENDINGPAYMENT:
		return orderv1.OrderStatusPENDINGPAYMENT, nil
	case model.OrderStatusPAID:
		return orderv1.OrderStatusPAID, nil
	case model.OrderStatusCANCELLED:
		return orderv1.OrderStatusCANCELLED, nil
	}
	return "", errs.ErrOrderInvalidStatus
}
