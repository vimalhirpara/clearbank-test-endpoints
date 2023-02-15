package handler

import (
	"bytes"
	"clearbanktestendpoints/pkg/model"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

// var _requestModel model.RequestModel
var authProfile model.AuthProfile

const cbUrl string = "https://institution-api-sim.clearbank.co.uk/v1/test"

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

func V1Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	if r.ContentLength > 0 {
		json.NewEncoder(w).Encode(model.SetErrorModel(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Do not require data, please remove.", nil))
		return
	}

	// Initialize Auth Profile: Token, Private Key, Public Key
	authProfile := model.InitAuthProfile()

	// Get Header
	postmanToken := r.Header.Get("Postman-Token") // Get Unique ID / UUID

	// Make the GET request
	request, err := http.NewRequest("GET", cbUrl, nil)
	if err != nil {
		json.NewEncoder(w).Encode(model.SetErrorModel(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Error creating GET request.", err))
		return
	}

	request.Header.Add("X-Request-Id", postmanToken)
	request.Header.Add("Authorization", "Bearer "+authProfile.Token)

	// Send the request and get the response
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		json.NewEncoder(w).Encode(model.SetErrorModel(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Error sending GET request.", err))
		return
	}

	defer response.Body.Close()

	json.NewEncoder(w).Encode(model.SetResponseModel(response.Status, time.Now(), ""))
}

func V1Post(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	// Get Header
	postmanToken := r.Header.Get("Postman-Token") // Get Unique ID / UUID

	if r.Header.Get("Content-Type") != "application/json" {
		json.NewEncoder(w).Encode(model.SetErrorModel(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid content type.", nil))
		return
	}

	if r.ContentLength == 0 {
		json.NewEncoder(w).Encode(model.SetErrorModel(http.StatusNoContent, http.StatusText(http.StatusNoContent), "Please send some data.", nil))
		return
	}

	// Initialize Auth Profile: Token, Private Key, Public Key
	authProfile := model.InitAuthProfile()

	// Decode Body
	//_ = json.NewDecoder(r.Body).Decode(&_requestModel)

	//apiRequestText, err := json.Marshal(requestModel{MachineName: _requestMachineName, UserName: _requestUserName, TimeStamp: _requestTimeStamp})
	apiRequestText, err := json.Marshal(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode(model.SetErrorModel(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Error marshal request body.", err))
		return
	}

	privateKey, err := loadPrivateKey()
	if err != nil {
		json.NewEncoder(w).Encode(model.SetErrorModel(http.StatusNotFound, http.StatusText(http.StatusNotFound), "Error private key not found/load.", err))
		return
	}

	dgtalSignature, err := generate(apiRequestText, privateKey)
	if err != nil {
		json.NewEncoder(w).Encode(model.SetErrorModel(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Error generate digital signature.", err))
		return
	}

	requestBody := bytes.NewBuffer(apiRequestText)

	request, err := http.NewRequest("POST", cbUrl, requestBody)
	if err != nil {
		json.NewEncoder(w).Encode(model.SetErrorModel(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Error creating POST request.", err))
		return
	}

	request.Header.Add("X-Request-Id", postmanToken)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("DigitalSignature", dgtalSignature)
	request.Header.Add("Authorization", "Bearer "+authProfile.Token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		json.NewEncoder(w).Encode(model.SetErrorModel(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Error sending POST request.", err))
		return
	}

	defer response.Body.Close()

	respBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		json.NewEncoder(w).Encode(model.SetErrorModel(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Error reading response.", err))
		return
	}
	_resp := model.ResponseModel{ResponseCode: response.Status, TimeStamp: time.Now(), Body: string(respBytes)}
	json.NewEncoder(w).Encode(_resp)
}

func generate(text []byte, privateKey *rsa.PrivateKey) (string, error) {
	rng := rand.Reader
	message := []byte(text)
	hashed := sha256.Sum256(message)

	signature, err := rsa.SignPKCS1v15(rng, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		//fmt.Fprintf(os.Stderr, "Error from signing: %s\n", err)
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func loadPrivateKey() (*rsa.PrivateKey, error) {
	priv, err := ioutil.ReadFile(authProfile.PrivateKeyPath)
	if err != nil {
		return nil, errors.New("no RSA private key found")
	}

	privPem, _ := pem.Decode(priv)
	if privPem.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("RSA private key is of the wrong type, Pem Type:" + privPem.Type)
	}
	privPemBytes := privPem.Bytes

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(privPemBytes); err != nil { // note this returns type `interface{}`
			return nil, errors.New("unable to parse RSA private key")
		}
	}

	var privateKey *rsa.PrivateKey
	var ok bool
	privateKey, ok = parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("unable to parse RSA private key")
	}

	return privateKey, nil
}
