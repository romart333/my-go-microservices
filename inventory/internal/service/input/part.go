package input

import "github.com/romart333/my-go-microservices/inventory/internal/model"

type PartFilter struct {
	// UUIDs — если не пустой, возвращаются только эти детали (приоритет)
	UUIDs []string
	// PartType — фильтр по типу (игнорируется если UUIDs заполнен)
	PartType model.PartType
}
