package tests

import (
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	api "github.com/romart333/my-go-microservices/inventory/internal/api/inventory/v1"
	"github.com/romart333/my-go-microservices/inventory/internal/api/inventory/v1/mocks"
	"github.com/romart333/my-go-microservices/inventory/internal/model"
	inventoryv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/inventory/v1"
)

func TestGetInventory(t *testing.T) {
	t.Parallel()

	type args struct {
		req *inventoryv1.GetPartRequest
	}

	type expected struct {
		err  error
		part *inventoryv1.Part
	}

	var (
		errService = errors.New("ошибка сервиса")

		inventoryUUID = uuid.MustParse(gofakeit.UUID())
		createdAt     = gofakeit.Date()
		price         = gofakeit.Int64()
		name          = gofakeit.Name()
		description   = gofakeit.Sentence()
		stockQuantity = gofakeit.Int64()

		part = model.Part{
			UUID:          inventoryUUID,
			Name:          name,
			Description:   description,
			Price:         price,
			PartType:      model.PartTypeHull,
			StockQuantity: stockQuantity,
			CreatedAt:     createdAt,
		}

		expectedPart = &inventoryv1.Part{
			Uuid:          inventoryUUID.String(),
			Name:          name,
			Description:   description,
			Price:         price,
			PartType:      inventoryv1.PartType_PART_TYPE_HULL,
			StockQuantity: stockQuantity,
			CreatedAt:     timestamppb.New(createdAt),
		}
	)

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(svc *mocks.PartService)
	}{
		{
			name: "успешное получение детали",
			args: args{
				req: &inventoryv1.GetPartRequest{Uuid: inventoryUUID.String()},
			},
			setupMock: func(svc *mocks.PartService) {
				svc.EXPECT().Get(t.Context(), inventoryUUID).Return(part, nil)
			},
			expected: expected{err: nil, part: expectedPart},
		},
		{
			name: "ошибка при получении детали",
			args: args{
				req: &inventoryv1.GetPartRequest{Uuid: inventoryUUID.String()},
			},
			setupMock: func(svc *mocks.PartService) {
				svc.EXPECT().Get(t.Context(), inventoryUUID).Return(model.Part{}, errService)
			},
			expected: expected{err: errService},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := mocks.NewPartService(t)
			if tt.setupMock != nil {
				tt.setupMock(svc)
			}
			server := api.NewInventoryServer(svc)
			resp, err := server.GetPart(t.Context(), tt.args.req)

			if tt.expected.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expected.err)
			} else {
				require.NoError(t, err)
				require.True(t, proto.Equal(tt.expected.part, resp.GetPart()))
			}
		})
	}
}
