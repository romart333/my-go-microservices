package converter

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/romart333/my-go-microservices/order/internal/model"
	paymentv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/payment/v1"
)

func ToPayInput(orderUUID uuid.UUID, method model.PaymentMethod) *paymentv1.PayOrderRequest {
	return &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		PaymentMethod: PaymentMethodToProto(method),
	}
}

func ToPayModel(response *paymentv1.PayOrderResponse) (uuid.UUID, error) {
	result, err := uuid.Parse(response.GetTransactionUuid())
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse transaction uuid: %w", err)
	}
	return result, nil
}

func PaymentMethodToProto(paymentMethod model.PaymentMethod) paymentv1.PaymentMethod {
	switch paymentMethod {
	case model.PaymentMethodCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSbp:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCreditCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodInvestorMoney:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
