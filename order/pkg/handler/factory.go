package handler

import (
	"errors"

	orderv1 "github.com/romart333/my-go-microservices/shared/pkg/openapi/order/v1"
	paymentv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/payment/v1"
)

func (h *OrderHandler) paymentMethodToProto(o orderv1.PaymentMethod) (paymentv1.PaymentMethod, error) {
	paymentMethodMap := map[orderv1.PaymentMethod]paymentv1.PaymentMethod{
		orderv1.PaymentMethodCARD:          paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
		orderv1.PaymentMethodSBP:           paymentv1.PaymentMethod_PAYMENT_METHOD_SBP,
		orderv1.PaymentMethodCREDITCARD:    paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		orderv1.PaymentMethodINVESTORMONEY: paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
	}
	protoPaymentMethod, ok := paymentMethodMap[o]
	if !ok {
		return protoPaymentMethod, errors.New("неверный payment_method")
	}

	return protoPaymentMethod, nil
}
