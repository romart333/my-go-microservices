package orderv1

import (
	"errors"
	"net/http"

	errs "github.com/romart333/my-go-microservices/order/internal/errors"
	orderv1 "github.com/romart333/my-go-microservices/shared/pkg/openapi/order/v1"
)

func handleCancelOrderError(err error) (orderv1.CancelOrderRes, error) {
	if errors.Is(err, errs.ErrOrderIdRequired) {
		return &orderv1.CancelOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "order_uuid обязательное поле",
		}, nil
	}
	if errors.Is(err, errs.ErrOrderCancelled) {
		return &orderv1.CancelOrderConflict{
			Code:    http.StatusConflict,
			Message: "Заказ уже отменён",
		}, nil
	}
	if errors.Is(err, errs.ErrOrderNotFound) {
		return &orderv1.CancelOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "Заказ не найден",
		}, nil
	}
	if errors.Is(err, errs.ErrOrderAlreadyPaid) {
		return &orderv1.CancelOrderConflict{
			Code:    http.StatusConflict,
			Message: "Заказ уже оплачен и не может быть отменён",
		}, nil
	}
	return &orderv1.CancelOrderInternalServerError{
		Code:    http.StatusInternalServerError,
		Message: "ошибка при отмене заказа",
	}, nil
}

func handleCreateOrderError(err error) (orderv1.CreateOrderRes, error) {
	if errors.Is(err, errs.ErrHullIsRequired) {
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "hull_uuid обязательное поле",
		}, nil
	}
	if errors.Is(err, errs.ErrEngineIsRequired) {
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "engine_uuid обязательное поле",
		}, nil
	}
	if errors.Is(err, errs.ErrOutOfStock) {
		return &orderv1.CreateOrderConflict{
			Code:    http.StatusConflict,
			Message: "деталей нет в наличии",
		}, nil
	}
	if errors.Is(err, errs.ErrPartNotFound) {
		return &orderv1.CreateOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "деталь не найдена",
		}, nil
	}
	return &orderv1.CreateOrderInternalServerError{
		Code:    http.StatusInternalServerError,
		Message: "ошибка при создании заказа",
	}, nil
}

func handleGetOrderError(err error) (orderv1.GetOrderRes, error) {
	if errors.Is(err, errs.ErrOrderNotFound) {
		return &orderv1.GetOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}
	if errors.Is(err, errs.ErrOrderIdRequired) {
		return &orderv1.GetOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "order_uuid обязательное поле",
		}, nil
	}
	return &orderv1.GetOrderInternalServerError{
		Code:    http.StatusInternalServerError,
		Message: "ошибка при получении заказа",
	}, nil
}

func handlePayOrderError(err error) (orderv1.PayOrderRes, error) {
	if errors.Is(err, errs.ErrOrderIdRequired) {
		return &orderv1.PayOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "order_uuid обязательное поле",
		}, nil
	}
	if errors.Is(err, errs.ErrOrderNotFound) {
		return &orderv1.PayOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}
	if errors.Is(err, errs.ErrInvalidParameters) {
		return &orderv1.PayOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "неверные параметры",
		}, nil
	}
	if errors.Is(err, errs.ErrPaymentNotFound) {
		return &orderv1.PayOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "информация по оплате заказа не найдена",
		}, nil
	}
	if errors.Is(err, errs.ErrOrderCancelled) {
		return &orderv1.PayOrderConflict{
			Code:    http.StatusConflict,
			Message: "заказ уже отменён",
		}, nil
	}
	if errors.Is(err, errs.ErrOrderAlreadyPaid) {
		return &orderv1.PayOrderConflict{
			Code:    http.StatusConflict,
			Message: "заказ уже оплачен",
		}, nil
	}
	return &orderv1.PayOrderInternalServerError{
		Code:    http.StatusInternalServerError,
		Message: "ошибка при оплате заказа",
	}, nil
}
