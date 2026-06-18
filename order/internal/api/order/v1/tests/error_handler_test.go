package tests

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	orderv1 "github.com/romart333/my-go-microservices/order/internal/api/order/v1"
	errs "github.com/romart333/my-go-microservices/order/internal/errors"
)

func TestErrorHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		err      error
		expected int
	}{
		// not found
		{
			name:     "заказ не найден",
			err:      errs.ErrOrderNotFound,
			expected: http.StatusNotFound,
		},
		{
			name:     "деталь не найдена",
			err:      errs.ErrPartNotFound,
			expected: http.StatusNotFound,
		},
		{
			name:     "информация по оплате заказа не найдена",
			err:      errs.ErrPaymentNotFound,
			expected: http.StatusNotFound,
		},
		// conflict
		{
			name:     "деталь отсутствует на складе",
			err:      errs.ErrOutOfStock,
			expected: http.StatusConflict,
		},
		{
			name:     "заказ уже оплачен",
			err:      errs.ErrOrderAlreadyPaid,
			expected: http.StatusConflict,
		},
		{
			name:     "заказ уже отменен",
			err:      errs.ErrOrderCancelled,
			expected: http.StatusConflict,
		},
		// bad request
		{
			name:     "неверный формат UUID",
			err:      errs.ErrInvalidUUID,
			expected: http.StatusBadRequest,
		},
		{
			name:     "неверные параметры",
			err:      errs.ErrInvalidParameters,
			expected: http.StatusBadRequest,
		},
		{
			name:     "неверный метод оплаты",
			err:      errs.ErrPaymentInvalidMethod,
			expected: http.StatusBadRequest,
		},
		{
			name:     "неизвестный тип детали",
			err:      errs.ErrPartUnknownType,
			expected: http.StatusBadRequest,
		},
		{
			name:     "неизвестный статус заказа",
			err:      errs.ErrOrderUnknownStatus,
			expected: http.StatusBadRequest,
		},
		// internal server error
		{
			name:     "внутренняя ошибка",
			err:      errors.New("внутренняя ошибка"),
			expected: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rr := httptest.NewRecorder()
			orderv1.ErrorHandler(t.Context(), rr, req, tt.err)

			assert.Equal(t, tt.expected, rr.Code)
		})
	}
}
