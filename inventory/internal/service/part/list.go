package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/romart333/my-go-microservices/inventory/internal/errors"
	"github.com/romart333/my-go-microservices/inventory/internal/model"
	"github.com/romart333/my-go-microservices/inventory/internal/service/input"
)

func (s *PartService) List(ctx context.Context, filter input.PartFilter) ([]model.Part, error) {
	for _, id := range filter.UUIDs {
		if id == "" {
			return nil, errs.ErrInvalidUUID
		}
		if _, err := uuid.Parse(id); err != nil {
			return nil, errs.ErrInvalidUUID
		}
	}
	parts, err := s.partRepository.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}
	return parts, nil
}
