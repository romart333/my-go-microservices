package v1

import (
	"context"

	"github.com/romart333/my-go-microservices/payment/internal/model"
)

type PaymentService interface {
	Pay(ctx context.Context, req model.PayRequest) (string, error)
}
