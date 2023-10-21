package controllers

type ResponseCode int

const (
	SuccessCode ResponseCode = 1000 + iota
	InvalidParamCode
	UserExistCode
	UserNotExistCode
	InvalidPasswordCode
)

var codeMsgMap = map[ResponseCode]string{
	SuccessCode:         "success",
	InvalidParamCode:    "请求参数有误",
	UserExistCode:       "用户名已存在",
	UserNotExistCode:    "用户名不存在",
	InvalidPasswordCode: "用户名或密码错误",
}

func (r ResponseCode) GetMsg(code ResponseCode) string {
	msg, ok := codeMsgMap[code]
	if !ok {
		return ""
	} else {
		return msg
	}
}
