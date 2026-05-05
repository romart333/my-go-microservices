package v1

import (
	"context"

	"github.com/romart333/my-go-microservices/inventory/internal/model"
)

type PartService interface {
	Get(ctx context.Context, id string) (model.Part, error)
	List(ctx context.Context, filter model.PartFilter) ([]model.Part, error)
}
