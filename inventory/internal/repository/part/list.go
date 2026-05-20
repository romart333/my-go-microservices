package part

import (
	"context"
	"fmt"
	"slices"
	"strings"

	errs "github.com/romart333/my-go-microservices/inventory/internal/errors"
	"github.com/romart333/my-go-microservices/inventory/internal/model"
	"github.com/romart333/my-go-microservices/inventory/internal/repository/converter"
	"github.com/romart333/my-go-microservices/inventory/internal/service/input"
)

func (r *repository) List(ctx context.Context, filter input.PartFilter) ([]model.Part, error) {
	r.RLock()
	defer r.RUnlock()
	parts := make([]model.Part, 0, len(r.parts))

	if len(filter.UUIDs) > 0 {
		for _, uuidStr := range filter.UUIDs {
			part, ok := r.parts[uuidStr]
			if !ok {
				return nil, fmt.Errorf("деталь с uuid %s не найдена: %w", uuidStr, errs.ErrPartNotFound)
			}
			parts = append(parts, converter.PartToModel(part))
		}
		return parts, nil
	}

	for _, part := range r.parts {
		partType := converter.PartTypeToModel(part.PartType)
		if filter.PartType == model.PartTypeUnspecified ||
			filter.PartType == partType {
			parts = append(parts, converter.PartToModel(part))
		}
	}
	slices.SortFunc(parts, func(a, b model.Part) int {
		return strings.Compare(a.Name, b.Name)
	})
	return parts, nil
}
