package v1

import (
	paymentv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/payment/v1"
)

type server struct {
	paymentv1.UnimplementedPaymentServiceServer
	paymentService PaymentService
}

func NewPaymentServer(paymentService PaymentService) *server {
	return &server{
		paymentService: paymentService,
	}
}
