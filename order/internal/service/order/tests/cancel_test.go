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

func TestCancel(t *testing.T) {
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
	)

	tests := []struct {
		name      string
		args      args
		setupMock func(repo *mocks.OrderRepository)
		expected  expected
	}{
		{
			name: "успешная отмена заказа",
			args: args{
				orderUUID: orderUUID,
			},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{
						UUID:   orderUUID,
						Status: model.OrderStatusPENDINGPAYMENT,
					}, nil)
				repo.EXPECT().
					Update(ctx, model.Order{
						UUID:   orderUUID,
						Status: model.OrderStatusCANCELLED,
					}).
					Return(nil)
			},
		},
		{
			name: "заказ уже оплачен",
			args: args{
				orderUUID: orderUUID,
			},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{
						UUID:   orderUUID,
						Status: model.OrderStatusPAID,
					}, nil)
			},
			expected: expected{err: errs.ErrOrderAlreadyPaid},
		},
		{
			name: "заказ уже отменен",
			args: args{
				orderUUID: orderUUID,
			},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().
					Get(ctx, orderUUID).
					Return(model.Order{
						UUID:   orderUUID,
						Status: model.OrderStatusCANCELLED,
					}, nil)
			},
			expected: expected{err: errs.ErrOrderCancelled},
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

			err := svc.Cancel(ctx, tc.args.orderUUID)
			if tc.expected.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expected.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
