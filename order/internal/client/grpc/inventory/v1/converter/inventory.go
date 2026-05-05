package converter

import (
	"github.com/romart333/my-go-microservices/order/internal/model"
	inventoryv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/inventory/v1"
)

func PartsToModel(parts []*inventoryv1.Part) []model.Part {
	modelParts := make([]model.Part, 0, len(parts))
	for _, part := range parts {
		modelParts = append(modelParts, model.Part{
			UUID:          part.GetUuid(),
			Name:          part.GetName(),
			Description:   part.GetDescription(),
			Price:         part.GetPrice(),
			PartType:      protoPartTypeToModel(part.GetPartType()),
			StockQuantity: part.GetStockQuantity(),
			CreatedAt:     new(part.GetCreatedAt().AsTime()),
		})
	}
	return modelParts
}

func protoPartTypeToModel(partType inventoryv1.PartType) model.PartType {
	switch partType {
	case inventoryv1.PartType_PART_TYPE_HULL:
		return model.PartTypeHull
	case inventoryv1.PartType_PART_TYPE_ENGINE:
		return model.PartTypeEngine
	case inventoryv1.PartType_PART_TYPE_SHIELD:
		return model.PartTypeShield
	default:
		return model.PartTypeUnspecified
	}
}
