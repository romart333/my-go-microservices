package handler

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	orderv1 "github.com/romart333/my-go-microservices/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/romart333/my-go-microservices/shared/pkg/proto/payment/v1"
)

const (
	middlewareTimeout = 10 * time.Second
	grpcCallTimeout   = 5 * time.Second
)

const (
	OrderStatusPENDINGPAYMENT OrderStatus = "PENDING_PAYMENT"
	OrderStatusPAID           OrderStatus = "PAID"
	OrderStatusCANCELLED      OrderStatus = "CANCELLED"
)

// OrderHandler реализует интерфейс orderv1.Handler, сгенерированный ogen.
type OrderHandler struct {
	orderv1.UnimplementedHandler
	inventoryClient inventoryv1.InventoryServiceClient
	paymentClient   paymentv1.PaymentServiceClient
	store           *OrderStore
}

// OrderStore — хранилище заказов (in-memory).
type OrderStore struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]Order
}

// Order представляет заказ на постройку космического корабля.
type Order struct {
	OrderUUID       uuid.UUID
	HullUUID        uuid.UUID
	EngineUUID      uuid.UUID
	ShieldUUID      *uuid.UUID // опциональный
	WeaponUUID      *uuid.UUID // опциональный
	TotalPrice      int64      // в копейках
	TransactionUUID *uuid.UUID
	PaymentMethod   *string
	Status          OrderStatus // PENDING_PAYMENT, PAID, CANCELLED
	CreatedAt       time.Time
}

type OrderStatus string

// NewOrderStore создаёт новое пустое хранилище заказов.
func NewOrderStore() *OrderStore {
	return &OrderStore{
		orders: make(map[uuid.UUID]Order),
	}
}

// NewOrderHandler создаёт новый обработчик заказов.
func NewOrderHandler(
	inventoryClient inventoryv1.InventoryServiceClient,
	paymentClient paymentv1.PaymentServiceClient,
	store *OrderStore,
) *OrderHandler {
	return &OrderHandler{
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		store:           store,
	}
}

// SetupServer создаёт OpenAPI сервер на основе обработчика.
func SetupServer(h *OrderHandler) (*orderv1.Server, error) {
	server, err := orderv1.NewServer(h)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания сервера OpenAPI: %w", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(middlewareTimeout))
	r.Use(middleware.Compress(5))
	r.Handle("/api/*", server)
	return server, nil
}

// GetOrder реализует операцию getOrder (пример реализации).
// GET /api/v1/orders/{order_uuid}.
func (h *OrderHandler) GetOrder(_ context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	// 1. Найти заказ в store (с блокировкой для thread-safety)
	h.store.mu.RLock()
	order, ok := h.store.orders[params.OrderUUID]
	h.store.mu.RUnlock()

	// 2. Если не найден — вернуть 404
	if !ok {
		return &orderv1.GetOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	// 3. Преобразовать в DTO и вернуть
	var shieldUUID orderv1.OptNilUUID
	if order.ShieldUUID != nil {
		shieldUUID = orderv1.NewOptNilUUID(*order.ShieldUUID)
	}

	var weaponUUID orderv1.OptNilUUID
	if order.WeaponUUID != nil {
		weaponUUID = orderv1.NewOptNilUUID(*order.WeaponUUID)
	}

	var transactionUUID orderv1.OptNilUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderv1.NewOptNilUUID(*order.TransactionUUID)
	}

	var paymentMethod orderv1.OptNilPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderv1.NewOptNilPaymentMethod(orderv1.PaymentMethod(*order.PaymentMethod))
	}

	return &orderv1.OrderDto{
		OrderUUID:       order.OrderUUID,
		HullUUID:        order.HullUUID,
		EngineUUID:      order.EngineUUID,
		ShieldUUID:      shieldUUID,
		WeaponUUID:      weaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderv1.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}, nil
}

// CreateOrder реализует операцию createOrder
// POST /api/v1/orders
func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	if req.GetHullUUID().String() == "" {
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "hull_uuid обязательное поле",
		}, nil
	}
	if req.GetEngineUUID().String() == "" {
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "engine_uuid обязательное поле",
		}, nil
	}

	uuids := []string{
		req.GetEngineUUID().String(),
		req.GetHullUUID().String(),
	}
	if req.GetShieldUUID().IsSet() {
		uuids = append(uuids, h.optNilUUIDToString(req.GetShieldUUID()))
	}
	if req.GetWeaponUUID().IsSet() {
		uuids = append(uuids, h.optNilUUIDToString(req.GetWeaponUUID()))
	}

	grpcCtx, cancel := context.WithTimeout(ctx, grpcCallTimeout)
	defer cancel()
	partsResponse, err := h.inventoryClient.ListParts(grpcCtx, &inventoryv1.ListPartsRequest{
		Uuids: uuids,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return &orderv1.CreateOrderNotFound{
					Code:    http.StatusNotFound,
					Message: "детали не найдены",
				}, nil
			case codes.InvalidArgument:
				return &orderv1.CreateOrderBadRequest{
					Code:    http.StatusBadRequest,
					Message: "неверный запрос",
				}, nil
			}
		}
		return &orderv1.CreateOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "внутренняя ошибка сервера при получении списка деталей",
		}, nil
	}

	totalPrice := int64(0)
	parts := partsResponse.GetParts()
	for _, part := range parts {
		if part.StockQuantity <= 0 {
			return &orderv1.CreateOrderConflict{
				Code:    http.StatusConflict,
				Message: fmt.Sprintf("детали %s нет в наличии", part.GetName()),
			}, nil
		}
		totalPrice += part.GetPrice()
	}

	orderId := uuid.New()

	var shieldUUID uuid.UUID
	if !req.GetShieldUUID().IsNull() {
		shieldUUID = req.GetShieldUUID().Value
	}
	var weaponUUID uuid.UUID
	if !req.GetWeaponUUID().IsNull() {
		weaponUUID = req.GetWeaponUUID().Value
	}
	order := Order{
		OrderUUID:  orderId,
		TotalPrice: totalPrice,
		Status:     OrderStatusPENDINGPAYMENT,
		CreatedAt:  time.Now(),
		HullUUID:   req.GetHullUUID(),
		EngineUUID: req.GetEngineUUID(),
		ShieldUUID: &shieldUUID,
		WeaponUUID: &weaponUUID,
	}
	h.store.mu.Lock()
	h.store.orders[orderId] = order
	h.store.mu.Unlock()

	return &orderv1.CreateOrderResponse{
		OrderUUID:  orderId,
		TotalPrice: totalPrice,
	}, nil
}

// PayOrder реализует операцию payOrder
// POST /api/v1/orders/{order_uuid}/pay
func (h *OrderHandler) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	h.store.mu.RLock()
	order, ok := h.store.orders[params.OrderUUID]
	h.store.mu.RUnlock()

	if !ok {
		return &orderv1.PayOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}
	if order.Status == OrderStatusCANCELLED {
		return &orderv1.PayOrderConflict{
			Code:    http.StatusConflict,
			Message: "Заказ уже отменён",
		}, nil
	}
	if order.Status == OrderStatusPAID {
		return &orderv1.PayOrderConflict{
			Code:    http.StatusConflict,
			Message: "Заказ уже оплачен",
		}, nil
	}

	protoPaymentMethod, err := h.paymentMethodToProto(req.GetPaymentMethod())
	if err != nil {
		return &orderv1.PayOrderBadRequest{
			Code:    http.StatusInternalServerError,
			Message: "внутренняя ошибка сервера",
		}, nil
	}

	grpcCtx, cancel := context.WithTimeout(ctx, grpcCallTimeout)
	defer cancel()
	response, err := h.paymentClient.PayOrder(grpcCtx, &paymentv1.PayOrderRequest{
		OrderUuid:     params.OrderUUID.String(),
		PaymentMethod: protoPaymentMethod,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return &orderv1.PayOrderBadRequest{
					Code:    http.StatusBadRequest,
					Message: "неверный запрос",
				}, nil
			case codes.NotFound:
				return &orderv1.PayOrderNotFound{
					Code:    http.StatusNotFound,
					Message: "информация по оплате заказа не найдена",
				}, nil
			}
		}
		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "внутренняя ошибка сервера при оплате заказа",
		}, nil
	}

	transactionUUID, err := uuid.Parse(response.GetTransactionUuid())
	if err != nil {
		return &orderv1.PayOrderBadRequest{
			Code:    http.StatusInternalServerError,
			Message: "внутренняя ошибка сервера при оплате заказа",
		}, nil
	}

	pm := string(req.GetPaymentMethod())
	h.store.mu.Lock()
	order.Status = OrderStatusPAID
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = &pm
	h.store.orders[params.OrderUUID] = order
	h.store.mu.Unlock()

	return &orderv1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}

// CancelOrder реализует операцию cancelOrder
// POST /api/v1/orders/{order_uuid}/cancel
func (h *OrderHandler) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	if params.OrderUUID.String() == "" {
		return &orderv1.CancelOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "order_uuid обязательное поле",
		}, nil
	}

	parsedUUID, err := uuid.Parse(params.OrderUUID.String())
	if err != nil {
		return &orderv1.CancelOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "order_uuid неверный UUID",
		}, nil
	}

	h.store.mu.RLock()
	order, ok := h.store.orders[parsedUUID]
	h.store.mu.RUnlock()
	if !ok {
		return &orderv1.CancelOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	switch order.Status {
	case OrderStatusPENDINGPAYMENT:
		order.Status = OrderStatusCANCELLED
		h.store.mu.Lock()
		h.store.orders[parsedUUID] = order
		h.store.mu.Unlock()
		return &orderv1.CancelOrderResponse{}, nil
	case OrderStatusPAID:
		return &orderv1.CancelOrderConflict{
			Code:    http.StatusConflict,
			Message: "Заказ уже оплачен и не может быть отменён",
		}, nil
	case OrderStatusCANCELLED:
		return &orderv1.CancelOrderConflict{
			Code:    http.StatusConflict,
			Message: "Заказ уже отменён",
		}, nil
	}

	return &orderv1.CancelOrderBadRequest{
		Code:    http.StatusBadRequest,
		Message: "неизвестный статус заказа",
	}, nil
}

func (h *OrderHandler) optNilUUIDToString(o orderv1.OptNilUUID) string {
	if v, ok := o.Get(); ok {
		return v.String()
	}
	return uuid.Nil.String()
}
