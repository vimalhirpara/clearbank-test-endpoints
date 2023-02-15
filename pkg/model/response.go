package model

import "time"

type ResponseModel struct {
	ResponseCode string    `json:"response-code"`
	TimeStamp    time.Time `json:"time-stamp,omitempty"`
	Body         string    `json:"body,omitempty"`
}

func SetResponseModel(responseCode string, timeStamp time.Time, body string) ResponseModel {
	return ResponseModel{ResponseCode: responseCode, TimeStamp: timeStamp, Body: body}
}
