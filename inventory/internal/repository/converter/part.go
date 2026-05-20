package converter

import (
	"github.com/romart333/my-go-microservices/inventory/internal/model"
	"github.com/romart333/my-go-microservices/inventory/internal/repository/record"
)

func PartToRecord(p model.Part) record.Part {
	return record.Part{
		UUID:          p.UUID,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		PartType:      partTypeToRecord(p.PartType),
		StockQuantity: p.StockQuantity,
		CreatedAt:     p.CreatedAt,
	}
}

func PartToModel(p record.Part) model.Part {
	return model.Part{
		UUID:          p.UUID,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		PartType:      PartTypeToModel(p.PartType),
		StockQuantity: p.StockQuantity,
		CreatedAt:     p.CreatedAt,
	}
}

func PartTypeToModel(p string) model.PartType {
	switch p {
	case "HULL":
		return model.PartTypeHull
	case "ENGINE":
		return model.PartTypeEngine
	case "SHIELD":
		return model.PartTypeShield
	case "WEAPON":
		return model.PartTypeWeapon
	}
	return model.PartTypeUnspecified
}

func partTypeToRecord(p model.PartType) string {
	switch p {
	case model.PartTypeHull:
		return "HULL"
	case model.PartTypeEngine:
		return "ENGINE"
	case model.PartTypeShield:
		return "SHIELD"
	case model.PartTypeWeapon:
		return "WEAPON"
	}
	return "UNSPECIFIED"
}
