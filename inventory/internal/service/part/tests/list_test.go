package tests

import (
	"context"
	"testing"

	gofakeit "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	errs "github.com/romart333/my-go-microservices/inventory/internal/errors"
	"github.com/romart333/my-go-microservices/inventory/internal/model"
	"github.com/romart333/my-go-microservices/inventory/internal/service/input"
	partservice "github.com/romart333/my-go-microservices/inventory/internal/service/part"
	"github.com/romart333/my-go-microservices/inventory/internal/service/part/mocks"
)

func TestList(t *testing.T) {
	t.Parallel()

	type args struct {
		filter input.PartFilter
	}

	type expected struct {
		parts []model.Part
		err   error
	}

	var (
		ctx = context.Background()

		parts = []model.Part{
			{
				UUID:          gofakeit.UUID(),
				Name:          gofakeit.Name(),
				Description:   gofakeit.Sentence(),
				Price:         gofakeit.Int64(),
				PartType:      model.PartType(model.PartType(gofakeit.RandomString([]string{string(model.PartTypeHull), string(model.PartTypeEngine), string(model.PartTypeShield), string(model.PartTypeWeapon)}))),
				StockQuantity: gofakeit.Int64(),
				CreatedAt:     gofakeit.Date(),
			},
			{
				UUID:          gofakeit.UUID(),
				Name:          gofakeit.Name(),
				Description:   gofakeit.Sentence(),
				Price:         gofakeit.Int64(),
				PartType:      model.PartType(model.PartType(gofakeit.RandomString([]string{string(model.PartTypeHull), string(model.PartTypeEngine), string(model.PartTypeShield), string(model.PartTypeWeapon)}))),
				StockQuantity: gofakeit.Int64(),
				CreatedAt:     gofakeit.Date(),
			},
		}
		filter = input.PartFilter{
			UUIDs:    []string{gofakeit.UUID()},
			PartType: model.PartType(model.PartType(gofakeit.RandomString([]string{string(model.PartTypeHull), string(model.PartTypeEngine), string(model.PartTypeShield), string(model.PartTypeWeapon)}))),
		}
	)

	tests := []struct {
		name      string
		args      args
		expected  expected
		setupMock func(partService *mocks.PartRepository)
	}{
		{
			name: "успешное получение списка деталей",
			args: args{
				filter: filter,
			},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().List(ctx, filter).Return(parts, nil)
			},
			expected: expected{parts: parts, err: nil},
		},
		{
			name: "один из UUID не найден",
			args: args{
				filter: filter,
			},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().List(ctx, filter).Return(nil, errs.ErrPartNotFound)
			},
			expected: expected{err: errs.ErrPartNotFound, parts: []model.Part{}},
		},
		{
			name: "ошибка валидации uuid",
			args: args{
				filter: input.PartFilter{
					UUIDs: []string{"invalid-uuid"},
				},
			},
			expected: expected{err: errs.ErrInvalidUUID, parts: []model.Part{}},
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
			parts, err := partService.List(ctx, tc.args.filter)
			if err != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expected.err)
			} else {
				require.NoError(t, err)
				require.ElementsMatch(t, tc.expected.parts, parts)
			}
		})
	}
}
