package part

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/romart333/my-go-microservices/inventory/internal/errors"
	"github.com/romart333/my-go-microservices/inventory/internal/model"
	"github.com/romart333/my-go-microservices/inventory/internal/repository/converter"
)

func (r *repository) Get(ctx context.Context, id uuid.UUID) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	part, ok := r.parts[id]
	if !ok {
		return model.Part{}, errs.ErrPartNotFound
	}
	return converter.PartToModel(part), nil
}
