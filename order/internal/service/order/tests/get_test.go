package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	errs "github.com/romart333/my-go-microservices/order/internal/errors"
	"github.com/romart333/my-go-microservices/order/internal/model"
	"github.com/romart333/my-go-microservices/order/internal/service/order"
	"github.com/romart333/my-go-microservices/order/internal/service/order/mocks"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type args struct {
		orderUUID string
	}

	type expected struct {
		err error
	}

	var (
		ctx       = context.Background()
		orderUUID = gofakeit.UUID()

		expectedOrder = model.Order{
			UUID:       orderUUID,
			HullUUID:   gofakeit.UUID(),
			EngineUUID: gofakeit.UUID(),
			TotalPrice: gofakeit.Int64(),
			Status:     model.OrderStatusPENDINGPAYMENT,
			CreatedAt:  gofakeit.Date(),
		}
	)

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(repo *mocks.OrderRepository)
	}{
		{
			name: "успешное получение заказа",
			args: args{
				orderUUID: orderUUID,
			},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(expectedOrder, nil)
			},
			expected: expected{err: nil},
		},
		{
			name: "заказ не найден",
			args: args{
				orderUUID: orderUUID,
			},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{}, errs.ErrOrderNotFound)
			},
			expected: expected{err: errs.ErrOrderNotFound},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repo := mocks.NewOrderRepository(t)
			inventoryClient := mocks.NewInventoryClient(t)
			paymentClient := mocks.NewPaymentClient(t)

			tc.setupMock(repo)

			svc := order.NewOrderService(inventoryClient, paymentClient, repo)

			order, err := svc.Get(ctx, tc.args.orderUUID)
			if tc.expected.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expected.err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, order.UUID)
				require.Equal(t, order, expectedOrder)
			}
		})
	}
}
