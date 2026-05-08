package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/romart333/my-go-microservices/order/internal/errors"
	"github.com/romart333/my-go-microservices/order/internal/model"
	"github.com/romart333/my-go-microservices/order/internal/service/order"
	"github.com/romart333/my-go-microservices/order/internal/service/order/mocks"
)

func TestPay(t *testing.T) {
	t.Parallel()

	type args struct {
		orderUUID       string
		method          model.PaymentMethod
		transactionUUID string
	}

	type expected struct {
		err             error
		transactionUUID string
	}

	var (
		ctx = context.Background()

		orderUUID = gofakeit.UUID()
		method    = model.PaymentMethod(gofakeit.RandomString([]string{
			string(model.PaymentMethodCARD),
			string(model.PaymentMethodSBP),
			string(model.PaymentMethodCREDITCARD),
			string(model.PaymentMethodINVESTORMONEY),
		}))
		transactionUUID = gofakeit.UUID()
	)

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(repo *mocks.OrderRepository, paymentClient *mocks.PaymentClient)
	}{
		{
			name: "успешная оплата заказа",
			args: args{
				orderUUID:       orderUUID,
				method:          method,
				transactionUUID: transactionUUID,
			},

			setupMock: func(repo *mocks.OrderRepository, paymentClient *mocks.PaymentClient) {
				repo.EXPECT().Get(ctx, orderUUID).Return(model.Order{
					UUID:            orderUUID,
					Status:          model.OrderStatusPENDINGPAYMENT,
					TransactionUUID: &transactionUUID,
					PaymentMethod:   &method,
				}, nil)
				paymentClient.EXPECT().PayOrder(ctx, orderUUID, method).Return(transactionUUID, nil)
				repo.EXPECT().Update(ctx, mock.MatchedBy(func(order model.Order) bool {
					return order.UUID == orderUUID &&
						order.Status == model.OrderStatusPAID &&
						*order.TransactionUUID == transactionUUID &&
						*order.PaymentMethod == method
				})).Return(nil)
			},
			expected: expected{
				transactionUUID: transactionUUID,
				err:             nil,
			},
		},
		{
			name: "заказ не найден",
			args: args{
				orderUUID:       orderUUID,
				method:          method,
				transactionUUID: transactionUUID,
			},
			setupMock: func(repo *mocks.OrderRepository, paymentClient *mocks.PaymentClient) {
				repo.EXPECT().Get(ctx, orderUUID).Return(model.Order{}, errs.ErrOrderNotFound)
			},
			expected: expected{
				err: errs.ErrOrderNotFound,
			},
		},
		{
			name: "заказ уже оплачен",
			args: args{
				orderUUID:       orderUUID,
				method:          method,
				transactionUUID: transactionUUID,
			},
			setupMock: func(repo *mocks.OrderRepository, paymentClient *mocks.PaymentClient) {
				repo.EXPECT().Get(ctx, orderUUID).Return(model.Order{
					UUID:   orderUUID,
					Status: model.OrderStatusPAID,
				}, nil)
			},
			expected: expected{
				err: errs.ErrOrderAlreadyPaid,
			},
		},
		{
			name: "заказ уже отменен",
			args: args{
				orderUUID:       orderUUID,
				method:          method,
				transactionUUID: transactionUUID,
			},
			setupMock: func(repo *mocks.OrderRepository, paymentClient *mocks.PaymentClient) {
				repo.EXPECT().Get(ctx, orderUUID).Return(model.Order{
					UUID:   orderUUID,
					Status: model.OrderStatusCANCELLED,
				}, nil)
			},
			expected: expected{
				err: errs.ErrOrderCancelled,
			},
		},
		{
			name: "ошибка PaymentService",
			args: args{
				orderUUID:       orderUUID,
				method:          method,
				transactionUUID: transactionUUID,
			},
			setupMock: func(repo *mocks.OrderRepository, paymentClient *mocks.PaymentClient) {
				repo.EXPECT().Get(ctx, orderUUID).Return(model.Order{
					UUID:   orderUUID,
					Status: model.OrderStatusPENDINGPAYMENT,
				}, nil)
				paymentClient.EXPECT().PayOrder(ctx, orderUUID, method).Return("", status.Error(codes.Internal, "internal"))
			},
			expected: expected{
				err: status.Error(codes.Internal, "internal"),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repo := mocks.NewOrderRepository(t)
			inventoryClient := mocks.NewInventoryClient(t)
			paymentClient := mocks.NewPaymentClient(t)

			tc.setupMock(repo, paymentClient)

			svc := order.NewOrderService(inventoryClient, paymentClient, repo)

			transactionUUID, err := svc.Pay(ctx, tc.args.orderUUID, tc.args.method)
			if tc.expected.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expected.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, transactionUUID, tc.expected.transactionUUID)
			}
		})
	}
}
