package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	errs "github.com/romart333/my-go-microservices/payment/internal/errors"
	model "github.com/romart333/my-go-microservices/payment/internal/model"
	input "github.com/romart333/my-go-microservices/payment/internal/service/input"
	service "github.com/romart333/my-go-microservices/payment/internal/service/payment"
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
		orderUUID = uuid.MustParse(gofakeit.UUID())
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
