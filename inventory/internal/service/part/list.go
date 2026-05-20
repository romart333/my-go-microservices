package part

import (
	"context"
	"fmt"

	"github.com/romart333/my-go-microservices/inventory/internal/model"
	"github.com/romart333/my-go-microservices/inventory/internal/service/input"
)

func (s *service) List(ctx context.Context, filter input.PartFilter) ([]model.Part, error) {
	parts, err := s.partRepository.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}
	return parts, nil
}
