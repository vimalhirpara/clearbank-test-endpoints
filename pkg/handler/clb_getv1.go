package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/vimalhirpara/clearbank-test-endpoints/pkg/model"
)

func V1Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	if r.ContentLength > 0 {
		json.NewEncoder(w).Encode(model.SetErrorModel(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "do not require data, please remove.", nil))
		return
	}

	// Initialize Auth Profile: Token, Private Key, Public Key
	authProfile := model.InitAuthProfile()

	// Get Header
	postmanToken := r.Header.Get("Postman-Token") // Get Unique ID / UUID

	// Make the GET request
	request, err := http.NewRequest("GET", cbUrl, nil)
	if err != nil {
		json.NewEncoder(w).Encode(model.SetErrorModel(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "creating GET request.", err))
		return
	}

	request.Header.Add("X-Request-Id", postmanToken)
	request.Header.Add("Authorization", "Bearer "+authProfile.Token)

	// Send the request and get the response
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		json.NewEncoder(w).Encode(model.SetErrorModel(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "sending GET request.", err))
		return
	}

	defer response.Body.Close()

	json.NewEncoder(w).Encode(model.SetResponseModel(response.Status, time.Now(), ""))
}
