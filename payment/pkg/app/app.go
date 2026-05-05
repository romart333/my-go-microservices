package app

import (
	"google.golang.org/grpc"

	api "github.com/romart333/my-go-microservices/payment/internal/api/payment/v1"
	"github.com/romart333/my-go-microservices/payment/internal/interceptor"
	service "github.com/romart333/my-go-microservices/payment/internal/service"
	paymentv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/payment/v1"
)

func RegisterServices(grpcServer *grpc.Server) {
	paymentService := service.NewPaymentService()
	paymentServer := api.NewPaymentServer(paymentService)
	paymentv1.RegisterPaymentServiceServer(grpcServer, paymentServer)
}

func Interceptors() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.LoggerInterceptor),
	}
}
