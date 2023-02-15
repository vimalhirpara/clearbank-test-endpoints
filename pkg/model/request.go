package model

import "time"

type RequestModel struct {
	MachineName string
	UserName    string
	TimeStamp   time.Time
}

func GetRequestModel() RequestModel {
	return RequestModel{MachineName: "CON-IND-LPT47", UserName: "Vimal Hirapara", TimeStamp: time.Now()}
}
