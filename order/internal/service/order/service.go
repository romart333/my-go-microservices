package order

type OrderService struct {
	inventoryClient InventoryClient
	paymentClient   PaymentClient
	orderRepository OrderRepository
}

func NewOrderService(inventoryClient InventoryClient, paymentClient PaymentClient, orderRepository OrderRepository) *OrderService {
	return &OrderService{
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		orderRepository: orderRepository,
	}
}
