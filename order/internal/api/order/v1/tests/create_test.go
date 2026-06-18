package tests

import (
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "github.com/romart333/my-go-microservices/order/internal/api/order/v1"
	"github.com/romart333/my-go-microservices/order/internal/api/order/v1/mocks"
	"github.com/romart333/my-go-microservices/order/internal/model"
	"github.com/romart333/my-go-microservices/order/internal/service/input"
	orderv1 "github.com/romart333/my-go-microservices/shared/pkg/openapi/order/v1"
)

func TestCreateOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		req *orderv1.CreateOrderRequest
	}

	type expected struct {
		err   error
		order *orderv1.CreateOrderResponse
	}

	var (
		errService = errors.New("ошибка сервиса")

		orderUUID   = uuid.MustParse(gofakeit.UUID())
		hullPrice   = gofakeit.Int64()
		enginePrice = gofakeit.Int64()

		hullUUID   = uuid.MustParse(gofakeit.UUID())
		engineUUID = uuid.MustParse(gofakeit.UUID())

		input = input.CreateOrderInput{
			HullUUID:   hullUUID,
			EngineUUID: engineUUID,
		}
	)
	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(orderService *mocks.OrderService)
	}{
		{
			name: "успешное создание заказа",
			args: args{
				req: &orderv1.CreateOrderRequest{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			expected: expected{
				err: nil,
				order: &orderv1.CreateOrderResponse{
					OrderUUID:  orderUUID,
					TotalPrice: hullPrice + enginePrice,
				},
			},
			setupMock: func(orderService *mocks.OrderService) {
				orderService.EXPECT().Create(t.Context(), input).Return(model.Order{
					UUID: orderUUID,
					Items: []model.OrderItem{
						{PartUUID: hullUUID, PartType: model.PartTypeHull, Price: hullPrice},
						{PartUUID: engineUUID, PartType: model.PartTypeEngine, Price: enginePrice},
					},
				}, nil)
			},
		},
		{
			name: "ошибка сервиса",
			args: args{
				req: &orderv1.CreateOrderRequest{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			expected: expected{err: errService},
			setupMock: func(orderService *mocks.OrderService) {
				orderService.EXPECT().Create(t.Context(), input).Return(model.Order{}, errService)
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

			actual, err := handler.CreateOrder(t.Context(), tt.args.req)

			if tt.expected.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expected.err)
			} else {
				actualOrder, ok := actual.(*orderv1.CreateOrderResponse)
				if !ok {
					t.Fatalf("expected *orderv1.CreateOrderResponse, got %T", actual)
				}
				require.NoError(t, err)
				assert.Equal(t, tt.expected.order.TotalPrice, actualOrder.TotalPrice)
				assert.Equal(t, tt.expected.order.OrderUUID, actualOrder.OrderUUID)
			}
		})
	}
}
