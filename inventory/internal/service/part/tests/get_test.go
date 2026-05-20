package tests

import (
	"context"
	"testing"

	gofakeit "github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	errs "github.com/romart333/my-go-microservices/inventory/internal/errors"
	"github.com/romart333/my-go-microservices/inventory/internal/model"
	partservice "github.com/romart333/my-go-microservices/inventory/internal/service/part"
	"github.com/romart333/my-go-microservices/inventory/internal/service/part/mocks"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type args struct {
		uuid uuid.UUID
	}

	type expected struct {
		part model.Part
		err  error
	}

	var (
		ctx = context.Background()

		uuid = uuid.MustParse(gofakeit.UUID())
		part = model.Part{
			UUID:          uuid,
			Name:          gofakeit.Name(),
			Description:   gofakeit.Sentence(),
			Price:         gofakeit.Int64(),
			PartType:      model.PartType(model.PartType(gofakeit.RandomString([]string{string(model.PartTypeHull), string(model.PartTypeEngine), string(model.PartTypeShield), string(model.PartTypeWeapon)}))),
			StockQuantity: gofakeit.Int64(),
			CreatedAt:     gofakeit.Date(),
		}
	)

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(partService *mocks.PartRepository)
	}{
		{
			name: "успешное получение детали",
			args: args{
				uuid: uuid,
			},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().Get(ctx, uuid).Return(part, nil)
			},
			expected: expected{part: part, err: nil},
		},
		{
			name: "деталь не найдена",
			args: args{
				uuid: uuid,
			},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().Get(ctx, uuid).Return(model.Part{}, errs.ErrPartNotFound)
			},
			expected: expected{err: errs.ErrPartNotFound, part: model.Part{}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewPartRepository(t)
			partService := partservice.NewPartService(repo)
			if tc.setupMock != nil {
				tc.setupMock(repo)
			}

			part, err := partService.Get(ctx, tc.args.uuid)
			if err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expected.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected.part, part)
			}
		})
	}
}
