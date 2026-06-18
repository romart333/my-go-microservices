package tests

import (
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	api "github.com/romart333/my-go-microservices/order/internal/api/order/v1"
	"github.com/romart333/my-go-microservices/order/internal/api/order/v1/mocks"
	orderv1 "github.com/romart333/my-go-microservices/shared/pkg/openapi/order/v1"
)

func TestCancelOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		params orderv1.CancelOrderParams
	}

	type expected struct {
		err   error
		order *orderv1.CancelOrderResponse
	}

	var (
		errService = errors.New("ошибка сервиса")

		orderUUID = uuid.MustParse(gofakeit.UUID())
	)

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(orderService *mocks.OrderService)
	}{
		{
			name: "успешная отмена заказа",
			args: args{
				params: orderv1.CancelOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{err: nil, order: &orderv1.CancelOrderResponse{}},
			setupMock: func(orderService *mocks.OrderService) {
				orderService.EXPECT().Cancel(t.Context(), orderUUID).Return(nil)
			},
		},
		{
			name: "ошибка сервиса",
			args: args{
				params: orderv1.CancelOrderParams{
					OrderUUID: orderUUID,
				},
			},
			expected: expected{err: errService},
			setupMock: func(orderService *mocks.OrderService) {
				orderService.EXPECT().Cancel(t.Context(), orderUUID).Return(errService)
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
			actual, err := handler.CancelOrder(t.Context(), tt.args.params)
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
