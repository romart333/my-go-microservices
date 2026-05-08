package v1

import (
	"context"

	"github.com/romart333/my-go-microservices/payment/internal/input"
)

type PaymentService interface {
	Pay(ctx context.Context, req input.PayOrderInput) (string, error)
}
