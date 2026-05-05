package v1

import (
	paymentv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/payment/v1"
)

// PaymentServer реализует gRPC сервис оплаты
type PaymentServer struct {
	paymentv1.UnimplementedPaymentServiceServer
	paymentService PaymentService
}

func NewPaymentServer(paymentService PaymentService) *PaymentServer {
	return &PaymentServer{
		paymentService: paymentService,
	}
}
