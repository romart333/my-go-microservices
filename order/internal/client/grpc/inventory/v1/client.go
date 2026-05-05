package v1

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/romart333/my-go-microservices/order/internal/client/grpc/inventory/v1/converter"
	errs "github.com/romart333/my-go-microservices/order/internal/errors"
	"github.com/romart333/my-go-microservices/order/internal/model"
	"github.com/romart333/my-go-microservices/shared/pkg/proto/inventory/v1"
)

const (
	grpcCallTimeout = 5 * time.Second
)

type Client struct {
	inventoryClient inventoryv1.InventoryServiceClient
}

func NewInventoryClient(inventoryClient inventoryv1.InventoryServiceClient) *Client {
	return &Client{inventoryClient: inventoryClient}
}

func (c *Client) ListParts(ctx context.Context, uuids []string) ([]model.Part, error) {
	grpcCtx, cancel := context.WithTimeout(ctx, grpcCallTimeout)
	defer cancel()
	resp, err := c.inventoryClient.ListParts(grpcCtx, &inventoryv1.ListPartsRequest{
		Uuids: uuids,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return nil, errs.ErrPartNotFound
		}
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}

	return converter.PartsToModel(resp.GetParts()), nil
}
