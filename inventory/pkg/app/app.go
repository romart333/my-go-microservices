package app

import (
	"google.golang.org/grpc"

	inventoryapi "github.com/romart333/my-go-microservices/inventory/internal/api/inventory/v1"
	"github.com/romart333/my-go-microservices/inventory/internal/interceptor"
	partrepo "github.com/romart333/my-go-microservices/inventory/internal/repository/part"
	partservice "github.com/romart333/my-go-microservices/inventory/internal/service/part"
	inventoryv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/inventory/v1"
)

func RegisterServices(grpcServer *grpc.Server) {
	partRepository := partrepo.NewPartRepository()
	partService := partservice.NewPartService(partRepository)
	inventoryServer := inventoryapi.NewInventoryServer(partService)

	inventoryv1.RegisterInventoryServiceServer(grpcServer, inventoryServer)
}

func Interceptors() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.ErrorInterceptor),
	}
}
