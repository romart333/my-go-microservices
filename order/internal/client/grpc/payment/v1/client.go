package payment

import (
	"context"
	"time"

	converter "github.com/romart333/my-go-microservices/order/internal/client/grpc/payment/v1/converter"
	"github.com/romart333/my-go-microservices/order/internal/model"
	paymentv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/payment/v1"
)

const (
	grpcCallTimeout = 5 * time.Second
)

type Client struct {
	paymentClient paymentv1.PaymentServiceClient
}

func NewPaymentClient(paymentClient paymentv1.PaymentServiceClient) *Client {
	return &Client{paymentClient: paymentClient}
}

func (c *Client) PayOrder(ctx context.Context, orderUUID string, method model.PaymentMethod) (string, error) {
	grpcCtx, cancel := context.WithTimeout(ctx, grpcCallTimeout)
	defer cancel()
	result, err := c.paymentClient.PayOrder(grpcCtx, &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID,
		PaymentMethod: converter.PaymentMethodToProto(method),
	})
	if err != nil {
		return "", err
	}
	transactionUUID := result.GetTransactionUuid()
	return transactionUUID, nil
}
