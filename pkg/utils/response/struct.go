package response

type SuccessResponseStruct struct {
	Code    int16       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Paging  interface{} `json:"paging,omitempty"`
	Extra   interface{} `json:"extra,omitempty"`
}

type ErrorResponseStruct struct {
	Code    int16       `json:"code"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}
