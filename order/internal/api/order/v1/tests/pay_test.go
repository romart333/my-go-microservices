package tests

import (
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	api "github.com/romart333/my-go-microservices/order/internal/api/order/v1"
	"github.com/romart333/my-go-microservices/order/internal/api/order/v1/mocks"
	"github.com/romart333/my-go-microservices/order/internal/model"
	orderv1 "github.com/romart333/my-go-microservices/shared/pkg/openapi/order/v1"
)

func TestPayOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		req    *orderv1.PayOrderRequest
		params orderv1.PayOrderParams
	}

	type expected struct {
		err   error
		order *orderv1.PayOrderResponse
	}

	var (
		errService = errors.New("ошибка сервиса")

		orderUUID       = uuid.MustParse(gofakeit.UUID())
		transactionUUID = uuid.MustParse(gofakeit.UUID())
	)
	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(orderService *mocks.OrderService)
	}{
		{
			name: "успешная оплата заказа",
			args: args{
				req: &orderv1.PayOrderRequest{
					PaymentMethod: orderv1.PaymentMethodCARD,
				},
				params: orderv1.PayOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{err: nil, order: &orderv1.PayOrderResponse{
				TransactionUUID: transactionUUID,
			}},
			setupMock: func(orderService *mocks.OrderService) {
				orderService.EXPECT().Pay(t.Context(), orderUUID, model.PaymentMethodCard).Return(transactionUUID, nil)
			},
		},
		{
			name: "ошибка сервиса",
			args: args{
				req: &orderv1.PayOrderRequest{
					PaymentMethod: orderv1.PaymentMethodCARD,
				},
				params: orderv1.PayOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{err: errService},
			setupMock: func(orderService *mocks.OrderService) {
				orderService.EXPECT().Pay(t.Context(), orderUUID, model.PaymentMethodCard).Return(uuid.Nil, errService)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			orderService := mocks.NewOrderService(t)
			handler := api.NewOrderHandler(orderService)
			if tt.setupMock != nil {
				tt.setupMock(orderService)
			}
			actual, err := handler.PayOrder(t.Context(), tt.args.req, tt.args.params)
			actualOrder, ok := actual.(*orderv1.PayOrderResponse)

			if tt.expected.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expected.err)
			} else {
				if !ok {
					t.Fatalf("expected *orderv1.PayOrderResponse, got %T", actual)
				}
				require.NoError(t, err)
				require.Equal(t, tt.expected.order, actualOrder)
			}
		})
	}
}
