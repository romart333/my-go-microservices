package v1

import (
	"context"

	"github.com/google/uuid"

	input "github.com/romart333/my-go-microservices/payment/internal/service/input"
)

type PaymentService interface {
	Pay(ctx context.Context, req input.PayOrderInput) (uuid.UUID, error)
}
