package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/vimalhirpara/clearbank-test-endpoints/pkg/model"
)

func V1Post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	// Get Header
	postmanToken := r.Header.Get("Postman-Token") // Get Unique ID / UUID

	if r.Header.Get("Content-Type") != "application/json" {
		json.NewEncoder(w).Encode(model.SetErrorModel(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "invalid content type.", nil))
		return
	}

	if r.ContentLength == 0 {
		json.NewEncoder(w).Encode(model.SetErrorModel(http.StatusNoContent, http.StatusText(http.StatusNoContent), "please send some data.", nil))
		return
	}

	// Initialize Auth Profile: Token, Private Key, Public Key
	authProfile := model.InitAuthProfile()

	// Decode Body
	//_ = json.NewDecoder(r.Body).Decode(&_requestModel)

	//apiRequestText, err := json.Marshal(requestModel{MachineName: _requestMachineName, UserName: _requestUserName, TimeStamp: _requestTimeStamp})
	apiRequestText, err := json.Marshal(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(model.SetErrorModel(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "marshal request body.", err))
		return
	}

	// Using HSM
	// digitalSignature, err := generateSignatureUsingHSM(string(apiRequestText))
	// if err != nil {
	// 	json.NewEncoder(w).Encode(model.SetErrorModel(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "generate digital signature.", err))
	// 	return
	// }

	privateKey, err := loadPrivateKey()
	if err != nil {
		json.NewEncoder(w).Encode(model.SetErrorModel(http.StatusNotFound, http.StatusText(http.StatusNotFound), "Error private key not found/load.", err))
		return
	}

	digitalSignature, err := generate(apiRequestText, privateKey)
	if err != nil {
		json.NewEncoder(w).Encode(model.SetErrorModel(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Error generate digital signature.", err))
		return
	}

	requestBody := bytes.NewBuffer(apiRequestText)

	request, err := http.NewRequest("POST", cbUrl, requestBody)
	if err != nil {
		json.NewEncoder(w).Encode(model.SetErrorModel(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "creating POST request.", err))
		return
	}

	request.Header.Add("X-Request-Id", postmanToken)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("DigitalSignature", digitalSignature)
	request.Header.Add("Authorization", "Bearer "+authProfile.Token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		json.NewEncoder(w).Encode(model.SetErrorModel(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "sending POST request.", err))
		return
	}

	defer response.Body.Close()

	respBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		json.NewEncoder(w).Encode(model.SetErrorModel(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "reading response.", err))
		return
	}
	_resp := model.ResponseModel{ResponseCode: response.Status, TimeStamp: time.Now(), Body: string(respBytes)}
	json.NewEncoder(w).Encode(_resp)
}
