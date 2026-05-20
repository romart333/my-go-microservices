package orderv1

import (
	"context"

	"github.com/google/uuid"

	"github.com/romart333/my-go-microservices/order/internal/model"
	"github.com/romart333/my-go-microservices/order/internal/service/input"
)

type OrderService interface {
	Create(ctx context.Context, req input.CreateOrderInput) (model.Order, error)
	Get(ctx context.Context, uuid uuid.UUID) (model.Order, error)
	Pay(ctx context.Context, uuid uuid.UUID, method model.PaymentMethod) (uuid.UUID, error)
	Cancel(ctx context.Context, uuid uuid.UUID) error
}
