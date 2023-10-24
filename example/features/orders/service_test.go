package orders

import (
	"fmt"
	"go-app/example/entity/message"
	"go-app/example/entity/model"
	"testing"
)

type orderRepositoryMock struct {
}

func (orderRepositoryMock) Save(value interface{}) error {
	order := value.(*model.Order)
	if order.OrderNumber == "1" {
		return fmt.Errorf("database error")
	}
	(*order).ID = 1
	return nil
}
func (orderRepositoryMock) SaveBatches(value interface{}, size int) error {
	return nil
}
func (orderRepositoryMock) DeleteOrderById(id int64) (int64, error) {
	return 1, nil
}
func (orderRepositoryMock) GetAllOrders() ([]model.Order, error) {
	return nil, nil
}

func TestServiceCreateOrderError(t *testing.T) {
	service := &orderService{
		repo: &orderRepositoryMock{},
	}
	_, err := service.CreateOrder(message.CreateOrder{
		OrderNumber: "1",
	})
	if err == nil {
		t.Error("it should return error")
	}
}
func TestServiceCreateOrderSuccess(t *testing.T) {
	service := &orderService{
		repo: &orderRepositoryMock{},
	}
	order, err := service.CreateOrder(message.CreateOrder{
		OrderNumber: "2",
	})
	if err != nil {
		t.Error("it should not return error")
	}
	if order.ID != 1 {
		t.Error("order id should be 1")
	}
}
