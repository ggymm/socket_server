package controllers

type ApiResponseJson struct {
	Status int         `json:"status"`
	Msg    interface{} `json:"msg"`
	Data   interface{} `json:"data"`
}

func ApiResource(status int, objects interface{}, msg string) (apiJson *ApiResponseJson) {
	apiJson = &ApiResponseJson{Status: status, Data: objects, Msg: msg}
	return
}

func ApiResourceSuccess(objects interface{}) (apiJson *ApiResponseJson) {
	apiJson = &ApiResponseJson{Status: 200, Data: objects, Msg: "success"}
	return
}

func ApiResourceError(msg string) (apiJson *ApiResponseJson) {
	apiJson = &ApiResponseJson{Status: 500, Data: nil, Msg: msg}
	return
}
