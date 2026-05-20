package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	errs "github.com/romart333/my-go-microservices/order/internal/errors"
	"github.com/romart333/my-go-microservices/order/internal/model"
	"github.com/romart333/my-go-microservices/order/internal/service/input"
)

func (s *service) Create(ctx context.Context, req input.CreateOrderInput) (model.Order, error) {
	uuids := req.PartUUIDs()
	parts, listErr := s.inventoryClient.ListParts(ctx, uuids)
	if listErr != nil {
		return model.Order{}, fmt.Errorf("получить список деталей: %w", listErr)
	}

	var orderItems []model.OrderItem
	for _, part := range parts {

		if part.StockQuantity <= 0 {
			return model.Order{}, fmt.Errorf("деталь %s: %w", part.Name, errs.ErrOutOfStock)
		}
		orderItems = append(orderItems, model.OrderItem{
			PartUUID: part.UUID,
			PartType: part.PartType,
			Price:    part.Price,
		})
	}

	orderId := uuid.New()

	order := model.Order{
		UUID:      orderId,
		Items:     orderItems,
		Status:    model.OrderStatusPendingPayment,
		CreatedAt: time.Now().UTC(),
	}
	createErr := s.orderRepository.Create(ctx, order)
	if createErr != nil {
		return model.Order{}, fmt.Errorf("создать заказ: %w", createErr)
	}

	return order, nil
}
