package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/romart333/my-go-microservices/inventory/internal/model"
	"github.com/romart333/my-go-microservices/inventory/internal/service/input"
)

type PartService interface {
	Get(ctx context.Context, id uuid.UUID) (model.Part, error)
	List(ctx context.Context, filter input.PartFilter) ([]model.Part, error)
}
