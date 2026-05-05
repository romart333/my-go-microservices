package orderv1

import (
	"context"

	"github.com/romart333/my-go-microservices/order/internal/model"
)

type OrderService interface {
	Create(ctx context.Context, req model.CreateOrderRequest) (model.Order, error)
	Get(ctx context.Context, uuid string) (model.Order, error)
	Pay(ctx context.Context, uuid string, method model.PaymentMethod) (string, error)
	Cancel(ctx context.Context, uuid string) error
}
