package main

import (
	"go-app/example/features"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Apis() http.Handler {
	router := httprouter.New()
	orderApi := features.GetOrderApi()
	router.GET("/orders", orderApi.ViewOrder)
	router.POST("/order", orderApi.CreateOrder)
	router.DELETE("/order/:id", orderApi.DeleteOrder)
	return router
}
