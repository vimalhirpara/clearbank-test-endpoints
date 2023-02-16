package hsmutil

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

var hsmClient HSMConnect

const (
	NoHash int = 0
	SHA1   int = 1
	MD5    int = 2
	RIPEMD int = 3
	SHA256 int = 4
)

const RSAS_MAX_LENGTH int = 4096 //Max length supported by HSM function

type HSMSigningResult struct {
	IsSuccess    bool
	IsContinue   bool
	SignedData   string
	ErrorMessage string
	CHCommand    string
}

type HSMSigning struct {
	primaryHSMIP   uint32
	secondaryHSMIP uint32
	tcpPort        int
	signKeyIndex   int16
}

func NewHSMSigningWithPort(primaryHSMLongIP uint32, secondaryHSMLongIP uint32, port int, signKeyIndex int16) *HSMSigning {
	var h HSMSigning
	h.primaryHSMIP = primaryHSMLongIP
	h.secondaryHSMIP = secondaryHSMLongIP
	h.tcpPort = port
	h.signKeyIndex = signKeyIndex

	fmt.Printf("Primary IP is: %d Secondary IP is: %d Port is: %d \n", h.primaryHSMIP, h.secondaryHSMIP, h.tcpPort)
	return &h
}

func (h *HSMSigning) GenerateSignature(inputMessage string) (string, error) {
	hashOfInputData := hex.EncodeToString([]byte(inputMessage))
	fmt.Println("SignData is", hashOfInputData)

	h.connectClient()

	HSMSignature := h.generateSign(hashOfInputData, h.signKeyIndex, SHA256)

	body, err := hex.DecodeString(HSMSignature)
	if err != nil {
		return "", err
	}

	//signature := base64.RawStdEncoding.EncodeToString(body)
	signature := base64.StdEncoding.EncodeToString(body)
	// fmt.Println("Digital Signature is: ", signature)

	return signature, nil
}

func (h *HSMSigning) generateSign(hexString string, signingKeyIndex int16, algorithm int) string {
	signedData := ""
	chCommand := ""

	objResult := HSMSigningResult{}
	objResult.IsSuccess = false
	hexChunkList := split(hexString, RSAS_MAX_LENGTH)

	for i := 0; i < len(hexChunkList); i++ {

		var isContinue bool = false
		if len(hexChunkList)-1 > 0 {
			isContinue = true
		}

		objResult = h.getSignData(
			hexChunkList[i],
			signingKeyIndex,
			isContinue,
			algorithm,
			func() string {
				if chCommand != "" {
					return chCommand
				}
				return ""
			}(),
		)

		chCommand = objResult.CHCommand

		if objResult.IsSuccess {
			if !objResult.IsContinue {
				signedData = objResult.SignedData
				break
			}
		} else {
			panic(fmt.Sprintf("Unable to sign data - %s", objResult.ErrorMessage))
		}
	}

	return signedData
}

func split(str string, chunkSize int) []string {
	var listArray []string
	remaining := 0

	for i := 0; i < len(str)/chunkSize; i++ {
		listArray = append(listArray, str[i*chunkSize:(i+1)*chunkSize])
	}

	if remaining*chunkSize < len(str) {
		listArray = append(listArray, str[remaining*chunkSize:])
	}

	return listArray
}

func (h *HSMSigning) getSignData(hexData string, privateKeyIndex int16, isContinue bool, algorithm int, sChcommand string) HSMSigningResult {
	var objResult HSMSigningResult
	objResult.IsSuccess = false

	var commandParam string
	var cmd string

	if strings.TrimSpace(sChcommand) == "" {
		commandParam = "[AORSAS;RC%d;RF%s;RG%d;BN%s;KY1;ZA1;]" // *KY1(BER encoding of the HASH); *ZA1(Padding (Default))
		cmd = fmt.Sprintf(commandParam,
			privateKeyIndex, // %d private key index
			hexData,         // %s Data used to generate the signature
			algorithm,       // %d Hash algorithm
			map[bool]string{true: "1", false: "0"}[isContinue]) // %s Send Data in chunk or not
	} else {
		// This section of command need to build when we have to pass CH parameter.
		commandParam = "[AORSAS;CH%s;RC%d;RF%s;RG%d;BN%s;KY1;ZA1;]" // *KY1(BER encoding of the HASH); *ZA1(Padding (Default))
		cmd = fmt.Sprintf(commandParam,
			sChcommand,      // %s CH command for split data
			privateKeyIndex, // %d private key index
			hexData,         // %s Data used to generate the signature
			algorithm,       // %d Hash algorithm
			map[bool]string{true: "1", false: "0"}[isContinue]) // %s Send Data in chunk or not
	}

	var endChar string = "]"
	response := h.executeExcrypt(cmd, endChar, !isContinue)
	functionID := ""

	response = response[1 : len(response)-1]
	resultArray := strings.Split(response, ";")

	var message string

	for i := 0; i < len(resultArray)-1; i++ {
		str := resultArray[i][0:2]
		data := resultArray[i][2:]
		switch str {
		case "AO":
			functionID = strings.ToUpper(data)
		case "BB":
			message = data
		case "BN":
			if strings.ToUpper(data) == "CONTINUE" {
				objResult.IsSuccess = true
				objResult.IsContinue = true
			} else {
				objResult.IsContinue = false
			}
		case "RH":
			objResult.IsSuccess = true
			objResult.SignedData = data
		case "CH":
			objResult.CHCommand = data
		}
	}

	if functionID == "ERRO" {
		objResult.IsSuccess = false
		objResult.ErrorMessage = message
	} else if functionID != "RSAS" {
		objResult.IsSuccess = false
		objResult.ErrorMessage = message
	}

	return objResult
}

func (h *HSMSigning) connectClient() {
	// fmt.Println("IP address is:", hsmPrimaryIP)

	hsmClient.NewHSMConnectWithPort(h.primaryHSMIP, h.secondaryHSMIP, h.tcpPort)

	//Connect to HSM
	err := hsmClient.Connect()

	if err != nil {
		panic(fmt.Sprintf("Unable to connect to Primary HSM - %s", err.Error()))
		// fmt.Errorf("unable to connect to Primary HSM, %w", err)
	}
}

func (h *HSMSigning) executeExcrypt(request, endChar string, endConnection bool) string {
	// now := time.Now()
	if !hsmClient.IsConnected() {
		h.connectClient()
		if !hsmClient.IsConnected() {
			return fmt.Sprintf("unable to Connect Primary %s and Secondary HSM %s", hsmClient.PrimaryHSMIP(), hsmClient.SecondaryHSMIP())
		}
	}

	result, err := hsmClient.PostRequest(request, endChar)
	if err != nil {
		panic(fmt.Sprintf("Unable to post request to HSM - %s", err.Error()))
		// fmt.Errorf("unable to post request to HSM, %w", err)
	}

	if endConnection {
		hsmClient.Disconnect()
	}

	return result
}
