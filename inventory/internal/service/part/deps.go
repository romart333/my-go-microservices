package part

import (
	"context"

	"github.com/romart333/my-go-microservices/inventory/internal/model"
)

type PartRepository interface {
	Get(ctx context.Context, id string) (model.Part, error)
	List(ctx context.Context, filter model.PartFilter) ([]model.Part, error)
}
