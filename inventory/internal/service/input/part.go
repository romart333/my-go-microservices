package input

import (
	"github.com/google/uuid"

	"github.com/romart333/my-go-microservices/inventory/internal/model"
)

type PartFilter struct {
	// UUIDs — если не пустой, возвращаются только эти детали (приоритет)
	UUIDs uuid.UUIDs
	// PartType — фильтр по типу (игнорируется если UUIDs заполнен)
	PartType model.PartType
}
