package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	errs "github.com/romart333/my-go-microservices/order/internal/errors"
	"github.com/romart333/my-go-microservices/order/internal/model"
)

func (s *OrderService) Create(ctx context.Context, req model.CreateOrderRequest) (model.Order, error) {
	if req.HullUUID == "" {
		return model.Order{}, errs.ErrHullIsRequired
	}

	if req.EngineUUID == "" {
		return model.Order{}, errs.ErrEngineIsRequired
	}

	uuids := req.PartUUIDs()
	parts, listErr := s.inventoryClient.ListParts(ctx, uuids)
	if listErr != nil {
		return model.Order{}, fmt.Errorf("получить список деталей: %w", listErr)
	}
	totalPrice := int64(0)
	for _, part := range parts {
		if part.StockQuantity <= 0 {
			return model.Order{}, fmt.Errorf("деталь %s: %w", part.Name, errs.ErrOutOfStock)
		}
		totalPrice += part.Price
	}

	orderId := uuid.New().String()

	order := model.Order{
		UUID:       orderId,
		TotalPrice: totalPrice,
		Status:     model.OrderStatusPENDINGPAYMENT,
		CreatedAt:  time.Now().UTC(),
		HullUUID:   req.HullUUID,
		EngineUUID: req.EngineUUID,
		ShieldUUID: req.ShieldUUID,
		WeaponUUID: req.WeaponUUID,
	}
	createErr := s.orderRepository.CreateOrder(order)
	if createErr != nil {
		return model.Order{}, fmt.Errorf("создать заказ: %w", createErr)
	}

	return order, nil
}
