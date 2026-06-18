package errs

import "errors"

var (
	// Ошибки заказов
	ErrOrderNotFound      = errors.New("заказ не найден")
	ErrOrderAlreadyPaid   = errors.New("заказ уже оплачен")
	ErrOrderCancelled     = errors.New("заказ отменён")
	ErrOrderUnknownStatus = errors.New("неизвестный статус заказа")

	// Ошибки деталей
	ErrPartNotFound    = errors.New("деталь не найдена")
	ErrPartUnknownType = errors.New("неизвестный тип детали")
	ErrOutOfStock      = errors.New("деталь отсутствует на складе")

	// Ошибки валидации (на границе API/Client конвертеров)
	ErrInvalidUUID = errors.New("неверный формат UUID")

	// Ошибки оплаты (приходят из PaymentClient)
	ErrPaymentNotFound      = errors.New("информация по оплате заказа не найдена")
	ErrPaymentInvalidMethod = errors.New("неверный метод оплаты")
	ErrInvalidParameters    = errors.New("неверные параметры")
)
