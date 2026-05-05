package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/romart333/my-go-microservices/inventory/internal/errors"
	"github.com/romart333/my-go-microservices/inventory/internal/model"
)

func (s *PartService) Get(ctx context.Context, id string) (model.Part, error) {
	if id == "" {
		return model.Part{}, errs.ErrInvalidUUID
	}

	if _, err := uuid.Parse(id); err != nil {
		return model.Part{}, errs.ErrInvalidUUID
	}

	part, err := s.partRepository.Get(ctx, id)
	if err != nil {
		return model.Part{}, fmt.Errorf("получить деталь: %w", err)
	}

	return part, nil
}
