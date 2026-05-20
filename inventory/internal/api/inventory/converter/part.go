package converter

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	errs "github.com/romart333/my-go-microservices/inventory/internal/errors"
	"github.com/romart333/my-go-microservices/inventory/internal/model"
	"github.com/romart333/my-go-microservices/inventory/internal/service/input"
	inventoryv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/inventory/v1"
)

func PartToProto(part model.Part) *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid:          part.UUID.String(),
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      PartTypeToProto(part.PartType),
		StockQuantity: part.StockQuantity,
		CreatedAt:     timestamppb.New(part.CreatedAt),
	}
}

func PartTypeToProto(partType model.PartType) inventoryv1.PartType {
	switch partType {
	case model.PartTypeHull:
		return inventoryv1.PartType_PART_TYPE_HULL
	case model.PartTypeEngine:
		return inventoryv1.PartType_PART_TYPE_ENGINE
	case model.PartTypeShield:
		return inventoryv1.PartType_PART_TYPE_SHIELD
	case model.PartTypeWeapon:
		return inventoryv1.PartType_PART_TYPE_WEAPON
	}
	return inventoryv1.PartType_PART_TYPE_UNSPECIFIED
}

func PartTypeToModel(partType inventoryv1.PartType) model.PartType {
	switch partType {
	case inventoryv1.PartType_PART_TYPE_HULL:
		return model.PartTypeHull
	case inventoryv1.PartType_PART_TYPE_ENGINE:
		return model.PartTypeEngine
	case inventoryv1.PartType_PART_TYPE_SHIELD:
		return model.PartTypeShield
	case inventoryv1.PartType_PART_TYPE_WEAPON:
		return model.PartTypeWeapon
	}
	return model.PartTypeUnspecified
}

func PartsToProto(parts []model.Part) []*inventoryv1.Part {
	protoParts := make([]*inventoryv1.Part, 0, len(parts))
	for _, part := range parts {
		protoParts = append(protoParts, PartToProto(part))
	}
	return protoParts
}

func ToGetInput(rawUUID string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(rawUUID)
	if err != nil {
		return uuid.Nil, errs.ErrInvalidUUID
	}
	return parsed, nil
}

func ProtoToPartFilter(request *inventoryv1.ListPartsRequest) (input.PartFilter, error) {
	uuids := make([]uuid.UUID, 0, len(request.GetUuids()))
	for _, uuidStr := range request.GetUuids() {
		parsed, err := uuid.Parse(uuidStr)
		if err != nil {
			return input.PartFilter{}, errs.ErrInvalidUUID
		}
		uuids = append(uuids, parsed)
	}
	return input.PartFilter{
		PartType: PartTypeToModel(request.GetPartType()),
		UUIDs:    uuids,
	}, nil
}
