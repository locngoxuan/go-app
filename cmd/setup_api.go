package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Apis() http.Handler {
	router := httprouter.New()
	return router
}
