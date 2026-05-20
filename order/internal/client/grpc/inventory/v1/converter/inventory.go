package converter

import (
	"fmt"

	"github.com/google/uuid"

	errs "github.com/romart333/my-go-microservices/order/internal/errors"
	"github.com/romart333/my-go-microservices/order/internal/model"
	inventoryv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/inventory/v1"
)

func PartsToModel(parts []*inventoryv1.Part) ([]model.Part, error) {
	modelParts := make([]model.Part, 0, len(parts))

	for _, part := range parts {
		uuid, err := uuid.Parse(part.GetUuid())
		if err != nil {
			return nil, fmt.Errorf("деталь с uuid %s: %w", part.GetUuid(), errs.ErrInvalidUUID)
		}
		modelParts = append(modelParts, model.Part{
			UUID:          uuid,
			Name:          part.GetName(),
			Description:   part.GetDescription(),
			Price:         part.GetPrice(),
			PartType:      protoPartTypeToModel(part.GetPartType()),
			StockQuantity: part.GetStockQuantity(),
			CreatedAt:     new(part.GetCreatedAt().AsTime()),
		})
	}
	return modelParts, nil
}

func InputToProto(uuids uuid.UUIDs) *inventoryv1.ListPartsRequest {
	parsedUuids := make([]string, 0, len(uuids))
	for _, id := range uuids {
		parsed := id.String()
		parsedUuids = append(parsedUuids, parsed)
	}
	return &inventoryv1.ListPartsRequest{
		Uuids: parsedUuids,
	}
}

func protoPartTypeToModel(partType inventoryv1.PartType) model.PartType {
	switch partType {
	case inventoryv1.PartType_PART_TYPE_HULL:
		return model.PartTypeHull
	case inventoryv1.PartType_PART_TYPE_ENGINE:
		return model.PartTypeEngine
	case inventoryv1.PartType_PART_TYPE_SHIELD:
		return model.PartTypeShield
	case inventoryv1.PartType_PART_TYPE_WEAPON:
		return model.PartTypeWeapon
	default:
		return model.PartTypeUnspecified
	}
}
