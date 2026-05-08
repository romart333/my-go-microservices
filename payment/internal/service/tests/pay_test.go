package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	errs "github.com/romart333/my-go-microservices/payment/internal/errors"
	"github.com/romart333/my-go-microservices/payment/internal/input"
	"github.com/romart333/my-go-microservices/payment/internal/model"
	"github.com/romart333/my-go-microservices/payment/internal/service"
)

func TestPay(t *testing.T) {
	t.Parallel()

	type args struct {
		input input.PayOrderInput
	}

	type expected struct {
		err error
	}

	var (
		ctx       = context.Background()
		orderUUID = gofakeit.UUID()
	)

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "успешная оплата картой",
			args: args{
				input: input.PayOrderInput{
					OrderUUID:     orderUUID,
					PaymentMethod: model.PaymentMethodCard,
				},
			},
		},
		{
			name: "успешная оплата картой по СБП",
			args: args{
				input: input.PayOrderInput{
					OrderUUID:     orderUUID,
					PaymentMethod: model.PaymentMethodSBP,
				},
			},
		},
		{
			name: "пустой order_uuid",
			args: args{
				input: input.PayOrderInput{
					OrderUUID:     "invalid-uuid",
					PaymentMethod: model.PaymentMethodSBP,
				},
			},
			expected: expected{
				err: errs.ErrInvalidOrderUUID,
			},
		},
		{
			name: "PaymentMethod = UNSPECIFIED",
			args: args{
				input: input.PayOrderInput{
					OrderUUID:     orderUUID,
					PaymentMethod: model.PaymentMethodUnspecified,
				},
			},
			expected: expected{
				err: errs.ErrInvalidPaymentMethod,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			paymentService := service.NewPaymentService()
			transactionUUID, err := paymentService.Pay(ctx, tc.args.input)
			if tc.expected.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expected.err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, transactionUUID)
			}
		})
	}
}
