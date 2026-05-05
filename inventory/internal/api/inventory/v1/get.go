package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/romart333/my-go-microservices/inventory/internal/converter"
	errs "github.com/romart333/my-go-microservices/inventory/internal/errors"
	inventoryv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/inventory/v1"
)

// GetPart возвращает деталь по UUID
func (s *InventoryServer) GetPart(
	ctx context.Context,
	req *inventoryv1.GetPartRequest,
) (*inventoryv1.GetPartResponse, error) {
	part, err := s.partService.Get(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, errs.ErrPartNotFound) {
			return nil, status.Error(codes.NotFound, "деталь не найдена")
		}
		if errors.Is(err, errs.ErrInvalidUUID) {
			return nil, status.Error(codes.InvalidArgument, "неверный формат uuid")
		}
		return nil, status.Errorf(codes.Internal, "ошибка при получении детали")
	}
	return &inventoryv1.GetPartResponse{
		Part: converter.PartToProto(part),
	}, nil
}
