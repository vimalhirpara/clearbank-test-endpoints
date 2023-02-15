package model

type ErrorModel struct {
	Status      int    `json:"status,omitempty"`
	Code        string `json:"code,omitempty"`
	Error       string `json:"error,omitempty"`
	ErrorDetail string `json:"error-detail,omitempty"`
}

func SetErrorModel(httpStatus int, httpCode string, _error string, errorDetail error) ErrorModel {
	if errorDetail != nil {
		return ErrorModel{Status: httpStatus, Code: httpCode, Error: _error, ErrorDetail: string(errorDetail.Error())}
	} else {
		return ErrorModel{Status: httpStatus, Code: httpCode, Error: _error}
	}
}
