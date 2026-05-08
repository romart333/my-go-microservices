package converter

import (
	"github.com/romart333/my-go-microservices/payment/internal/input"
	"github.com/romart333/my-go-microservices/payment/internal/model"
	paymentv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/payment/v1"
)

func PayRequestToInput(req *paymentv1.PayOrderRequest) input.PayOrderInput {
	return input.PayOrderInput{
		OrderUUID:     req.GetOrderUuid(),
		PaymentMethod: PaymentMethodToModel(req.GetPaymentMethod()),
	}
}

func PaymentMethodToModel(pm paymentv1.PaymentMethod) model.PaymentMethod {
	switch pm {
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case paymentv1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnspecified
	}
}
