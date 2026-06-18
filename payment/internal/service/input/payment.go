package input

import (
	"github.com/google/uuid"

	"github.com/romart333/my-go-microservices/payment/internal/model"
)

type PayOrderInput struct {
	OrderUUID     uuid.UUID
	PaymentMethod model.PaymentMethod
}
