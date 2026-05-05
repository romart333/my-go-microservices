package v1

import (
	inventoryv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/inventory/v1"
)

type InventoryServer struct {
	inventoryv1.UnimplementedInventoryServiceServer
	partService PartService
}

func NewInventoryServer(partService PartService) *InventoryServer {
	return &InventoryServer{
		partService: partService,
	}
}
