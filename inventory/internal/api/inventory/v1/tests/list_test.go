package tests

import (
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	api "github.com/romart333/my-go-microservices/inventory/internal/api/inventory/v1"
	"github.com/romart333/my-go-microservices/inventory/internal/api/inventory/v1/mocks"
	"github.com/romart333/my-go-microservices/inventory/internal/model"
	"github.com/romart333/my-go-microservices/inventory/internal/service/input"
	inventoryv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/inventory/v1"
)

func TestListInventory(t *testing.T) {
	t.Parallel()

	type args struct {
		req *inventoryv1.ListPartsRequest
	}

	type expected struct {
		err   error
		parts *inventoryv1.ListPartsResponse
	}

	var (
		errService = errors.New("ошибка сервиса")
		uuid1      = uuid.MustParse(gofakeit.UUID())
		uuid2      = uuid.MustParse(gofakeit.UUID())

		createdAt     = gofakeit.Date()
		price         = gofakeit.Int64()
		stockQuantity = gofakeit.Int64()
		name          = gofakeit.Name()
		description   = gofakeit.Sentence()
		part1         = model.Part{
			UUID:          uuid1,
			Name:          name,
			Description:   description,
			Price:         price,
			PartType:      model.PartTypeHull,
			StockQuantity: stockQuantity,
			CreatedAt:     createdAt,
		}
		part2 = model.Part{
			UUID:          uuid2,
			Name:          name,
			Description:   description,
			Price:         price,
			PartType:      model.PartTypeHull,
			StockQuantity: stockQuantity,
			CreatedAt:     createdAt,
		}
		parts = []model.Part{
			part1,
			part2,
		}

		filter = input.PartFilter{
			PartType: model.PartTypeHull,
			UUIDs: []uuid.UUID{
				uuid1,
				uuid2,
			},
		}
	)

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(inventoryService *mocks.PartService)
	}{
		{
			name: "успешное получение списка деталей",
			args: args{
				req: &inventoryv1.ListPartsRequest{
					PartType: inventoryv1.PartType_PART_TYPE_HULL,
					Uuids: []string{
						uuid1.String(),
						uuid2.String(),
					},
				},
			},
			setupMock: func(inventoryService *mocks.PartService) {
				inventoryService.EXPECT().List(t.Context(), filter).Return(parts, nil)
			},
			expected: expected{
				err: nil,
				parts: &inventoryv1.ListPartsResponse{
					Parts: []*inventoryv1.Part{
						{
							Uuid:          uuid1.String(),
							Name:          name,
							Description:   description,
							Price:         price,
							PartType:      inventoryv1.PartType_PART_TYPE_HULL,
							StockQuantity: stockQuantity,
							CreatedAt:     timestamppb.New(createdAt),
						},
						{
							Uuid:          uuid2.String(),
							Name:          name,
							Description:   description,
							Price:         price,
							PartType:      inventoryv1.PartType_PART_TYPE_HULL,
							StockQuantity: stockQuantity,
							CreatedAt:     timestamppb.New(createdAt),
						},
					},
				},
			},
		},
		{
			name: "ошибка при получении списка деталей",
			args: args{
				req: &inventoryv1.ListPartsRequest{
					PartType: inventoryv1.PartType_PART_TYPE_HULL,
					Uuids: []string{
						uuid1.String(),
						uuid2.String(),
					},
				},
			},
			setupMock: func(inventoryService *mocks.PartService) {
				inventoryService.EXPECT().List(t.Context(), filter).Return(nil, errService)
			},
			expected: expected{
				err:   errService,
				parts: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			inventoryService := mocks.NewPartService(t)
			server := api.NewInventoryServer(inventoryService)
			tt.setupMock(inventoryService)

			actual, err := server.ListParts(t.Context(), tt.args.req)

			if tt.expected.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expected.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected.parts, actual)
			}
		})
	}
}
