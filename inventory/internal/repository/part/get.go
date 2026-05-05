package part

import (
	"context"

	errs "github.com/romart333/my-go-microservices/inventory/internal/errors"
	"github.com/romart333/my-go-microservices/inventory/internal/model"
	"github.com/romart333/my-go-microservices/inventory/internal/repository/converter"
)

func (r *PartRepository) Get(ctx context.Context, id string) (model.Part, error) {
	r.RLock()
	defer r.RUnlock()
	part, ok := r.parts[id]
	if !ok {
		return model.Part{}, errs.ErrPartNotFound
	}
	return converter.PartToModel(part), nil
}
