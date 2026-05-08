package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	errs "github.com/romart333/my-go-microservices/order/internal/errors"
	"github.com/romart333/my-go-microservices/order/internal/model"
	"github.com/romart333/my-go-microservices/order/internal/service/input"
	"github.com/romart333/my-go-microservices/order/internal/service/order"
	"github.com/romart333/my-go-microservices/order/internal/service/order/mocks"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type args struct {
		in input.CreateOrderInput
	}

	type expected struct {
		err error
	}

	var (
		ctx = context.Background()

		hullUUID   = gofakeit.UUID()
		engineUUID = gofakeit.UUID()

		partsInStock = []model.Part{
			{UUID: hullUUID, Name: "Hull", Price: 500000, StockQuantity: 10},
			{UUID: engineUUID, Name: "Engine", Price: 300000, StockQuantity: 5},
		}

		partsOutOfStock = []model.Part{
			{UUID: hullUUID, Name: "Hull", Price: 500000, StockQuantity: 10},
			{UUID: engineUUID, Name: "Engine", Price: 300000, StockQuantity: 0},
		}
	)

	tests := []struct {
		name      string
		args      args
		setupMock func(repo *mocks.OrderRepository, client *mocks.InventoryClient)
		expected  expected
	}{
		{
			name: "успешное создание заказа",
			args: args{
				in: input.CreateOrderInput{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().
					ListParts(ctx, []string{hullUUID, engineUUID}).
					Return(partsInStock, nil)

				repo.EXPECT().
					Create(ctx, mock.MatchedBy(func(o model.Order) bool {
						return o.HullUUID == hullUUID &&
							o.EngineUUID == engineUUID &&
							o.TotalPrice == 800000 && // 500000 + 300000
							o.Status == model.OrderStatusPENDINGPAYMENT
					})).
					Return(nil)
			},
			expected: expected{err: nil},
		},
		{
			name: "деталь не найдена",
			args: args{
				in: input.CreateOrderInput{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().
					ListParts(ctx, []string{hullUUID, engineUUID}).
					Return(nil, errs.ErrPartNotFound)
			},
			expected: expected{err: errs.ErrPartNotFound},
		},
		{
			name: "деталь закончилась на складе",
			args: args{
				in: input.CreateOrderInput{
					HullUUID:   hullUUID,
					EngineUUID: engineUUID,
				},
			},
			setupMock: func(repo *mocks.OrderRepository, client *mocks.InventoryClient) {
				client.EXPECT().
					ListParts(ctx, []string{hullUUID, engineUUID}).
					Return(partsOutOfStock, nil)
			},
			expected: expected{err: errs.ErrOutOfStock},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewOrderRepository(t)
			inventoryClient := mocks.NewInventoryClient(t)
			paymentClient := mocks.NewPaymentClient(t)

			tc.setupMock(repo, inventoryClient)

			svc := order.NewOrderService(inventoryClient, paymentClient, repo)

			order, err := svc.Create(ctx, tc.args.in)
			if tc.expected.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expected.err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, order.UUID)
			}
		})
	}
}
