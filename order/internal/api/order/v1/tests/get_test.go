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

func TestGetOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		orderUUID uuid.UUID
	}

	type expected struct {
		err      error
		orderDTO *orderv1.OrderDto
	}

	var (
		errService = errors.New("ошибка сервиса")

		orderUUID       = uuid.MustParse(gofakeit.UUID())
		transactionUUID = uuid.MustParse(gofakeit.UUID())

		hullUUID   = uuid.MustParse(gofakeit.UUID())
		engineUUID = uuid.MustParse(gofakeit.UUID())
		shieldUUID = uuid.MustParse(gofakeit.UUID())
		weaponUUID = uuid.MustParse(gofakeit.UUID())

		hullPrice   = gofakeit.Int64()
		enginePrice = gofakeit.Int64()
		shieldPrice = gofakeit.Int64()
		weaponPrice = gofakeit.Int64()

		paymentMethod = model.PaymentMethodCreditCard
		status        = model.OrderStatusPendingPayment
		createdAt     = gofakeit.Date()

		serviceOrder = model.Order{
			UUID: orderUUID,
			Items: []model.OrderItem{
				{PartUUID: hullUUID, PartType: model.PartTypeHull, Price: hullPrice},
				{PartUUID: engineUUID, PartType: model.PartTypeEngine, Price: enginePrice},
				{PartUUID: shieldUUID, PartType: model.PartTypeShield, Price: shieldPrice},
				{PartUUID: weaponUUID, PartType: model.PartTypeWeapon, Price: weaponPrice},
			},
			Status:          status,
			CreatedAt:       createdAt,
			PaymentMethod:   &paymentMethod,
			TransactionUUID: &transactionUUID,
		}

		orderDTO = &orderv1.OrderDto{
			OrderUUID:       orderUUID,
			TotalPrice:      serviceOrder.TotalPrice(),
			Status:          orderv1.OrderStatusPENDINGPAYMENT,
			CreatedAt:       createdAt,
			PaymentMethod:   orderv1.NewOptNilPaymentMethod(orderv1.PaymentMethodCREDITCARD),
			TransactionUUID: orderv1.NewOptNilUUID(transactionUUID),
			HullUUID:        hullUUID,
			EngineUUID:      engineUUID,
			ShieldUUID:      orderv1.NewOptNilUUID(shieldUUID),
			WeaponUUID:      orderv1.NewOptNilUUID(weaponUUID),
		}
	)
	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(orderService *mocks.OrderService)
	}{
		{
			name: "успешное получение заказа",
			args: args{
				orderUUID: orderUUID,
			},
			expected: expected{err: nil, orderDTO: orderDTO},
			setupMock: func(orderService *mocks.OrderService) {
				orderService.EXPECT().Get(t.Context(), orderUUID).Return(serviceOrder, nil)
			},
		},
		{
			name: "ошибка сервиса",
			args: args{
				orderUUID: orderUUID,
			},
			expected: expected{err: errService, orderDTO: nil},
			setupMock: func(orderService *mocks.OrderService) {
				orderService.EXPECT().Get(t.Context(), orderUUID).Return(model.Order{}, errService)
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

			actual, err := handler.GetOrder(t.Context(), orderv1.GetOrderParams{OrderUUID: tt.args.orderUUID})

			if tt.expected.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expected.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected.orderDTO, actual)
			}
		})
	}
}
