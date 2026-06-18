package order

type service struct {
	inventoryClient InventoryClient
	paymentClient   PaymentClient
	orderRepository OrderRepository
}

func NewOrderService(inventoryClient InventoryClient, paymentClient PaymentClient, orderRepository OrderRepository) *service {
	return &service{
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		orderRepository: orderRepository,
	}
}
