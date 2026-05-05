package errs

import "errors"

var (
	ErrOrderNotFound      = errors.New("заказ не найден")
	ErrOrderAlreadyPaid   = errors.New("заказ уже оплачен")
	ErrOrderCancelled     = errors.New("заказ отменён")
	ErrOrderUnknownStatus = errors.New("неизвестный статус заказа")
	ErrPartNotFound       = errors.New("деталь не найдена")
	ErrOutOfStock         = errors.New("деталь отсутствует на складе")
	ErrPaymentNotFound    = errors.New("информация по оплате заказа не найдена")
	ErrOrderIdRequired    = errors.New("order id обязательное поле")

	ErrHullIsRequired   = errors.New("корпус обязательное поле")
	ErrEngineIsRequired = errors.New("двигатель обязательное поле")

	ErrInvalidParameters = errors.New("неверные параметры")

	ErrOrderInvalidStatus        = errors.New("неверный статус заказа")
	ErrOrderInvalidPaymentMethod = errors.New("неверный метод оплаты")
)

// var (
//     // Ошибки заказов
//     ErrOrderNotFound    = errors.New("заказ не найден")
//     ErrOrderAlreadyPaid = errors.New("заказ уже оплачен")
//     ErrOrderCancelled   = errors.New("заказ отменён")
//     ErrOrderAlreadyCancelled = errors.New("заказ уже отменён")
//     ErrOrderUnknownStatus = errors.New("неизвестный статус заказа")

//     // Ошибки деталей
//     ErrPartNotFound = errors.New("деталь не найдена")
//     ErrOutOfStock   = errors.New("деталь отсутствует на складе")
//     ErrHullIsRequired = errors.New("корпус обязательное поле")
//     ErrEngineIsRequired = errors.New("двигатель обязательное поле")

//     // Ошибки валидации
//     ErrInvalidUUID          = errors.New("неверный формат UUID")
//     ErrInvalidPaymentMethod = errors.New("неверный метод оплаты")
//     ErrOrderIdRequired      = errors.New("order id обязательное поле")
// )
