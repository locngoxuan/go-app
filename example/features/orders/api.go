package orders

import (
	"encoding/json"
	"go-app/example/entity/message"
	"go-app/example/helper"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"go.uber.org/dig"
)

type OrderApi interface {
	CreateOrder(http.ResponseWriter, *http.Request, httprouter.Params)
	ViewOrder(http.ResponseWriter, *http.Request, httprouter.Params)
	DeleteOrder(http.ResponseWriter, *http.Request, httprouter.Params)
}

type orderApiParams struct {
	dig.In
	OrderService
}

type orderApi struct {
	orderService OrderService
}

func NewOrderApi(params orderApiParams) OrderApi {
	return &orderApi{
		orderService: params.OrderService,
	}
}

// CreateOrder implements OrderApi.
func (o *orderApi) CreateOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var request message.CreateOrder
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	if helper.IsAnyBlank(request.Customer, request.OrderNumber) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	order, err := o.orderService.CreateOrder(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	bytes, _ := json.Marshal(order)
	w.WriteHeader(200)
	w.Write(bytes)
}

// DeleteOrder implements OrderApi.
func (api *orderApi) DeleteOrder(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	err = api.orderService.CancelOrder(int64(id))
	if err != nil {
		log.Error().Err(err).Msg("failed to cancel order")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

// ViewOrder implements OrderApi.
func (api *orderApi) ViewOrder(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	orders, err := api.orderService.ViewOrder()
	if err != nil {
		log.Error().Err(err).Msg("faild to get order from database")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
	data, _ := json.Marshal(orders)
	w.WriteHeader(200)
	w.Write(data)
}
