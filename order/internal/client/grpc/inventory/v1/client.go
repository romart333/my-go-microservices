package v1

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/romart333/my-go-microservices/order/internal/client/grpc/inventory/v1/converter"
	errs "github.com/romart333/my-go-microservices/order/internal/errors"
	"github.com/romart333/my-go-microservices/order/internal/model"
	inventoryv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/inventory/v1"
)

const (
	grpcCallTimeout = 5 * time.Second
)

type client struct {
	inventoryClient inventoryv1.InventoryServiceClient
}

func NewInventoryClient(inventoryClient inventoryv1.InventoryServiceClient) *client {
	return &client{inventoryClient: inventoryClient}
}

func (c *client) ListParts(ctx context.Context, uuids uuid.UUIDs) ([]model.Part, error) {
	grpcCtx, cancel := context.WithTimeout(ctx, grpcCallTimeout)
	defer cancel()
	resp, err := c.inventoryClient.ListParts(grpcCtx, converter.InputToProto(uuids))
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return nil, errs.ErrPartNotFound
		}
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}

	parts, err := converter.PartsToModel(resp.GetParts())
	if err != nil {
		return nil, fmt.Errorf("конвертация деталей: %w", err)
	}
	return parts, nil
}
