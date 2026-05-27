package tests

import (
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	api "github.com/romart333/my-go-microservices/payment/internal/api/payment/v1"
	"github.com/romart333/my-go-microservices/payment/internal/api/payment/v1/mocks"
	"github.com/romart333/my-go-microservices/payment/internal/model"
	"github.com/romart333/my-go-microservices/payment/internal/service/input"
	paymentv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/payment/v1"
)

func TestPayOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		req *paymentv1.PayOrderRequest
	}

	type expected struct {
		err   error
		order *paymentv1.PayOrderResponse
	}

	var (
		errService = errors.New("ошибка сервиса")

		orderUUID       = uuid.MustParse(gofakeit.UUID())
		paymentMethod   = model.PaymentMethodCard
		transactionUUID = uuid.MustParse(gofakeit.UUID())
	)

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(paymentService *mocks.PaymentService)
	}{
		{
			name: "успешная оплата",
			args: args{
				req: &paymentv1.PayOrderRequest{
					OrderUuid:     orderUUID.String(),
					PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
				},
			},
			expected: expected{
				err: nil,
				order: &paymentv1.PayOrderResponse{
					TransactionUuid: transactionUUID.String(),
				},
			},
			setupMock: func(paymentService *mocks.PaymentService) {
				paymentService.EXPECT().Pay(t.Context(), input.PayOrderInput{
					OrderUUID:     orderUUID,
					PaymentMethod: paymentMethod,
				}).Return(transactionUUID, nil)
			},
		},
		{
			name: "ошибка сервиса",
			args: args{
				req: &paymentv1.PayOrderRequest{
					OrderUuid:     orderUUID.String(),
					PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
				},
			},
			expected: expected{
				err:   errService,
				order: nil,
			},
			setupMock: func(paymentService *mocks.PaymentService) {
				paymentService.EXPECT().Pay(t.Context(), input.PayOrderInput{
					OrderUUID:     orderUUID,
					PaymentMethod: paymentMethod,
				}).Return(uuid.UUID{}, errService)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			paymentService := mocks.NewPaymentService(t)
			server := api.NewPaymentServer(paymentService)
			tt.setupMock(paymentService)

			actual, err := server.PayOrder(t.Context(), tt.args.req)

			if tt.expected.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expected.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected.order, actual)
			}
		})
	}
}
