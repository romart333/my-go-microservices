package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/romart333/my-go-microservices/inventory/internal/model"
)

func (s *service) Get(ctx context.Context, id uuid.UUID) (model.Part, error) {
	part, err := s.partRepository.Get(ctx, id)
	if err != nil {
		return model.Part{}, fmt.Errorf("получить деталь: %w", err)
	}

	return part, nil
}
