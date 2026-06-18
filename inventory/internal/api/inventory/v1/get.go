package v1

import (
	"context"
	"fmt"
	"log/slog"

	converter "github.com/romart333/my-go-microservices/inventory/internal/api/inventory/converter"
	inventoryv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/inventory/v1"
)

// GetPart возвращает деталь по UUID
func (s *server) GetPart(
	ctx context.Context,
	req *inventoryv1.GetPartRequest,
) (*inventoryv1.GetPartResponse, error) {
	partUUID, err := converter.ToGetInput(req.GetUuid())
	if err != nil {
		slog.ErrorContext(ctx, "получить деталь: разбор UUID", "uuid", req.GetUuid(), "error", err)
		return nil, fmt.Errorf("разобрать UUID детали: %w", err)
	}

	part, err := s.partService.Get(ctx, partUUID)
	if err != nil {
		slog.ErrorContext(ctx, "получить деталь", "part_uuid", partUUID, "error", err)
		return nil, err
	}

	return &inventoryv1.GetPartResponse{
		Part: converter.PartToProto(part),
	}, nil
}
