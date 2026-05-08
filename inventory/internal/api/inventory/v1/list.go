package v1

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/romart333/my-go-microservices/inventory/internal/converter"
	errs "github.com/romart333/my-go-microservices/inventory/internal/errors"
	"github.com/romart333/my-go-microservices/inventory/internal/service/input"
	inventoryv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/inventory/v1"
)

// ListParts возвращает список деталей с опциональной фильтрацией по типу
func (s *InventoryServer) ListParts(
	ctx context.Context,
	req *inventoryv1.ListPartsRequest,
) (*inventoryv1.ListPartsResponse, error) {
	filter := input.PartFilter{
		PartType: converter.PartTypeToModel(req.GetPartType()),
		UUIDs:    req.GetUuids(),
	}
	parts, err := s.partService.List(ctx, filter)
	if err != nil {
		slog.ErrorContext(ctx, "получение списка деталей", "error", err)
		if errors.Is(err, errs.ErrInvalidUUID) {
			return nil, status.Errorf(codes.InvalidArgument, "неверный формат uuid")
		}
		if errors.Is(err, errs.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "деталь не найдена")
		}
		return nil, status.Errorf(codes.Internal, "ошибка при получении списка деталей")
	}
	return &inventoryv1.ListPartsResponse{
		Parts: converter.PartsToProto(parts),
	}, nil
}
