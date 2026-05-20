package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	converter "github.com/romart333/my-go-microservices/order/internal/client/grpc/payment/v1/converter"
	errs "github.com/romart333/my-go-microservices/order/internal/errors"
	"github.com/romart333/my-go-microservices/order/internal/model"
	paymentv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/payment/v1"
)

const (
	grpcCallTimeout = 5 * time.Second
)

type client struct {
	paymentClient paymentv1.PaymentServiceClient
}

func NewPaymentClient(paymentClient paymentv1.PaymentServiceClient) *client {
	return &client{paymentClient: paymentClient}
}

func (c *client) PayOrder(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	grpcCtx, cancel := context.WithTimeout(ctx, grpcCallTimeout)
	defer cancel()
	result, err := c.paymentClient.PayOrder(grpcCtx, converter.ToPayInput(orderUUID, method))
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return uuid.Nil, errs.ErrInvalidParameters
			case codes.NotFound:
				return uuid.Nil, errs.ErrPaymentNotFound
			}
		}
		return uuid.Nil, fmt.Errorf("оплатить заказ: %w", err)
	}

	transactionUUID, err := converter.ToPayModel(result)
	if err != nil {
		return uuid.Nil, fmt.Errorf("convert payment response to output: %w", err)
	}
	return transactionUUID, nil
}
