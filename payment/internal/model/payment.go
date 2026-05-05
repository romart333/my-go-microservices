package model

const (
	PaymentMethodUnspecified   PaymentMethod = "UNSPECIFIED"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

type PayRequest struct {
	OrderUUID     string
	PaymentMethod PaymentMethod
}

type PaymentMethod string

func (m PaymentMethod) IsValid() bool {
	switch m {
	case PaymentMethodCard, PaymentMethodSBP,
		PaymentMethodCreditCard, PaymentMethodInvestorMoney:
		return true
	default:
		return false
	}
}
