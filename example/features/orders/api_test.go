package orders

import (
	"bytes"
	"encoding/json"
	"go-app/example/entity/message"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type orderServiceMock struct {
}

func (*orderServiceMock) CreateOrder(request message.CreateOrder) (*message.Order, error) {
	if request.OrderNumber == "1" {
		return &message.Order{
			ID:          1,
			Customer:    request.Customer,
			OrderNumber: request.OrderNumber,
			Amount:      request.Amount,
			Created:     time.Now().Unix(),
		}, nil
	}
	return nil, nil
}

func (*orderServiceMock) CancelOrder(orderId int64) error {
	return nil
}
func (*orderServiceMock) ViewOrder() ([]message.Order, error) {
	return nil, nil
}

func TestApiCreateOrderReturnBadRequestForEmptyPayload(t *testing.T) {
	api := &orderApi{
		orderService: &orderServiceMock{},
	}
	req := httptest.NewRequest("POST", "/order", nil)
	w := httptest.NewRecorder()
	api.CreateOrder(w, req, nil)
	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("it should return 400")
	}
}

func TestApiCreateOrderReturnBadRequest(t *testing.T) {
	api := &orderApi{
		orderService: &orderServiceMock{},
	}
	payload, _ := json.Marshal(map[string]interface{}{})
	req := httptest.NewRequest("POST", "/order", bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	api.CreateOrder(w, req, nil)
	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("it should return 400")
	}
}
func TestApiCreateOrderReturnOK(t *testing.T) {
	api := &orderApi{
		orderService: &orderServiceMock{},
	}
	payload, _ := json.Marshal(map[string]interface{}{
		"customer":    "customer",
		"orderNumber": "1",
		"amount":      100,
	})
	req := httptest.NewRequest("POST", "/order", bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	api.CreateOrder(w, req, nil)
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("it should return 200")
	}
}

func TestApiCancelOrderReturnBadRequest(t *testing.T) {
	api := &orderApi{
		orderService: &orderServiceMock{},
	}
	req := httptest.NewRequest("DELETE", "/order/order-1", nil)
	w := httptest.NewRecorder()
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("it should return 200")
	}
	api.DeleteOrder(w, req, nil)
}

func TestApiCancelOrderReturnOK(t *testing.T) {
	api := &orderApi{
		orderService: &orderServiceMock{},
	}
	req := httptest.NewRequest("DELETE", "/order/1", nil)
	w := httptest.NewRecorder()
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("it should return 200")
	}
	api.DeleteOrder(w, req, nil)
}
