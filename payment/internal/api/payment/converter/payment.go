package converter

import (
	"github.com/google/uuid"

	errs "github.com/romart333/my-go-microservices/payment/internal/errors"
	"github.com/romart333/my-go-microservices/payment/internal/model"
	"github.com/romart333/my-go-microservices/payment/internal/service/input"
	paymentv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/payment/v1"
)

func PayRequestToInput(req *paymentv1.PayOrderRequest) (input.PayOrderInput, error) {
	orderUUID, err := uuid.Parse(req.GetOrderUuid())
	if err != nil {
		return input.PayOrderInput{}, errs.ErrInvalidOrderUUID
	}

	return input.PayOrderInput{
		OrderUUID:     orderUUID,
		PaymentMethod: PaymentMethodToModel(req.GetPaymentMethod()),
	}, nil
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
