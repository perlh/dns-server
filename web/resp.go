package web

import "encoding/json"

type Result[T interface{}] struct {
	Data  T                      `json:"data" desc:"data"`
	Code  int                    `json:"code" desc:"状态码"`
	State bool                   `json:"state" desc:"true,false"`
	Msg   string                 `json:"msg" desc:"提示信息"`
	Attr  map[string]interface{} `json:"attr" desc:"熟悉，其他附加参数"`
}

// UnSafeJsonSerializer 不安全的json序列化，并返回Result
func UnSafeJsonSerializer(data interface{}) []byte {
	res, err := json.Marshal(data)
	if err != nil {
		res, _ = json.Marshal(ResultFailureMsg[string](err.Error()))
	}
	return res
}

func ResultFailureMsg[T interface{}](msg string) *Result[T] {
	return &Result[T]{
		Code:  400,
		State: false,
		Msg:   msg,
	}
}

func ResultFailureCodeMsg[T interface{}](code int, msg string) *Result[T] {
	return &Result[T]{
		Code:  code,
		State: false,
		Msg:   msg,
	}
}

func ResultSuccessMsg[T interface{}](msg string) *Result[T] {
	return &Result[T]{
		Code:  200,
		State: true,
		Msg:   msg,
	}
}

func ResultSuccessMsgAttr[T interface{}](msg string, attr map[string]interface{}) *Result[T] {
	return &Result[T]{
		Code:  200,
		State: true,
		Msg:   msg,
		Attr:  attr,
	}
}

func ResultSuccessDataMsg[T interface{}](data T, msg string) *Result[T] {
	return &Result[T]{
		Data:  data,
		Code:  200,
		State: true,
		Msg:   msg,
	}
}

func ResultSuccessDataMsgAttr[T interface{}](data T, msg string, attr map[string]interface{}) *Result[T] {
	return &Result[T]{
		Data:  data,
		Code:  200,
		State: true,
		Msg:   msg,
		Attr:  attr,
	}
}

func ResultSuccessData[T interface{}](data T) *Result[T] {
	return &Result[T]{
		Data:  data,
		Code:  200,
		State: true,
		Msg:   "success",
	}
}

func ResultSuccessDataCode[T interface{}](data T, code int) *Result[T] {
	return &Result[T]{
		Data:  data,
		Code:  code,
		State: true,
		Msg:   "success",
	}
}

func ResultSuccessDataCodeMsg[T interface{}](data T, code int, msg string) *Result[T] {
	return &Result[T]{
		Data:  data,
		Code:  code,
		State: true,
		Msg:   msg,
	}
}

func ResultSuccessDataAttrMsg[T interface{}](data T, attr map[string]interface{}, msg string) *Result[T] {
	return &Result[T]{
		Data:  data,
		Code:  200,
		State: true,
		Attr:  attr,
		Msg:   msg,
	}
}

func ResultSuccessDataAttr[T interface{}](data T, attr map[string]interface{}) *Result[T] {
	return &Result[T]{
		Data:  data,
		Code:  200,
		State: true,
		Attr:  attr,
		Msg:   "success",
	}
}
