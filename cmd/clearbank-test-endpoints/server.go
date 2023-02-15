package main

import (
	"clearbanktestendpoints/pkg/handler"

	"github.com/gorilla/mux"
)

func initRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", handler.Welcome).Methods("GET")
	router.HandleFunc("/healthcheck", handler.Healthcheck).Methods("GET")
	router.HandleFunc("/v1get", handler.V1Get).Methods("GET")
	router.HandleFunc("/v1post", handler.V1Post).Methods("POST")
	return router
}
