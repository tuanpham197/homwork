package response

func SuccessResponse(code int16, mesage string, data, paging, extra interface{}) *SuccessResponseStruct {
	return &SuccessResponseStruct{Code: code, Message: mesage, Data: data, Paging: paging, Extra: extra}
}

func ErrorResponse(code int16, message string, err interface{}) *ErrorResponseStruct {
	return &ErrorResponseStruct{Code: code, Message: message, Errors: err}
}

func ResponseData(data interface{}, code int16, message string) *SuccessResponseStruct {
	return SuccessResponse(code, message, data, nil, nil)
}
