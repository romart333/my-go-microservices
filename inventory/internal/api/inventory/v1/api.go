package v1

import (
	inventoryv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/inventory/v1"
)

type server struct {
	inventoryv1.UnimplementedInventoryServiceServer
	partService PartService
}

func NewInventoryServer(partService PartService) *server {
	return &server{
		partService: partService,
	}
}
