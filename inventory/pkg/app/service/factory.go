package service

import (
	inventoryv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/inventory/v1"
)

func (s *InventoryServer) modelToProto(part Part) *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      part.PartType,
		StockQuantity: part.StockQuantity,
		CreatedAt:     part.CreatedAt,
	}
}
