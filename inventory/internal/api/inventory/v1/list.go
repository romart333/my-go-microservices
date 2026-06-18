package v1

import (
	"context"
	"fmt"
	"log/slog"

	converter "github.com/romart333/my-go-microservices/inventory/internal/api/inventory/converter"
	inventoryv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/inventory/v1"
)

// ListParts возвращает список деталей с опциональной фильтрацией по типу
func (s *server) ListParts(
	ctx context.Context,
	req *inventoryv1.ListPartsRequest,
) (*inventoryv1.ListPartsResponse, error) {
	filter, err := converter.ProtoToPartFilter(req)
	if err != nil {
		slog.ErrorContext(ctx, "получить список деталей: разбор фильтра", "error", err)
		return nil, fmt.Errorf("разобрать фильтр: %w", err)
	}

	parts, err := s.partService.List(ctx, filter)
	if err != nil {
		slog.ErrorContext(ctx, "получить список деталей", "error", err)
		return nil, err
	}

	return &inventoryv1.ListPartsResponse{
		Parts: converter.PartsToProto(parts),
	}, nil
}
