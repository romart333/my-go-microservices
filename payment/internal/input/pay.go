package input

import "github.com/romart333/my-go-microservices/payment/internal/model"

type PayOrderInput struct {
	OrderUUID     string
	PaymentMethod model.PaymentMethod
}
