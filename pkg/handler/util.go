package handler

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io/ioutil"

	"github.com/vimalhirpara/clearbank-test-endpoints/pkg/model"

	hsmutil "github.com/AartiChhasiya/generate-digital-signature/hsm"
)

const cbUrl string = "https://institution-api-sim.clearbank.co.uk/v1/test"

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

	// Initialize Auth Profile: Token, Private Key, Public Key
	authProfile := model.InitAuthProfile()

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

func generateSignatureUsingHSM(inputMessage string) (string, error) {
	var hsmPrimaryIP uint32 = 3232272392
	var hsmSecondaryIP uint32 = 3232272392
	var hsmPort int = 9000
	var signKeyIndex int16 = 23

	hsm := hsmutil.NewHSMSigningWithPort(hsmPrimaryIP, hsmSecondaryIP, hsmPort, signKeyIndex)
	signature, err := hsm.GenerateSignature(inputMessage)
	return signature, err
}
