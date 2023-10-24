package orders

import (
	"go-app/example/entity/message"
	"go-app/example/entity/model"
	"time"

	"go.uber.org/dig"
)

type OrderService interface {
	CreateOrder(request message.CreateOrder) (*model.Order, error)
	CancelOrder(orderId int64) error
	ViewOrder() ([]message.Order, error)
}

type orderServiceParams struct {
	dig.In
	OrderRepo OrderRepository
}

func NewOrderService(params orderServiceParams) OrderService {
	return &orderService{
		repo: params.OrderRepo,
	}
}

type orderService struct {
	repo OrderRepository
}

// CancelOrder implements OrderService.
func (service *orderService) CancelOrder(orderId int64) error {
	_, err := service.repo.DeleteOrderById(orderId)
	return err
}

// CreateOrder implements OrderService.
func (service *orderService) CreateOrder(request message.CreateOrder) (*model.Order, error) {
	order := model.Order{
		Customer:    request.Customer,
		OrderNumber: request.OrderNumber,
		Amount:      request.Amount,
		Created:     time.Now().Unix(),
	}
	err := service.repo.Save(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// ViewOrder implements OrderService.
func (service *orderService) ViewOrder() ([]message.Order, error) {
	orders, err := service.repo.GetAllOrders()
	if err != nil {
		return nil, err
	}
	result := make([]message.Order, 0)
	for _, order := range orders {
		result = append(result, message.Order{
			ID:          order.ID,
			Customer:    order.Customer,
			OrderNumber: order.OrderNumber,
			Amount:      order.Amount,
			Created:     order.Created,
		})
	}
	return result, nil
}
