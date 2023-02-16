package handler

import (
	"encoding/json"
	"net/http"

	"github.com/vimalhirpara/clearbank-test-endpoints/pkg/model"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/text")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")
	w.Write([]byte("Welcome to Clear Bank test environment."))
}

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")
	json.NewEncoder(w).Encode(model.GetRequestModel())
}
