package contextcode

import "github.com/gin-gonic/gin"

const (
	ContextUserIDKey   = "userID"
	ContextUsernameKey = "username"
	ExpireTimeKey      = "expireTime"
)

type Response struct {
	Code    ResponseCode `json:"code"`
	Message interface{}  `json:"message"`
	Data    interface{}  `json:"data"`
}

func ReturnResponse(c *gin.Context, httpStatus int, code ResponseCode, data ...interface{}) {
	r := Response{
		Code:    code,
		Message: code.GetMsg(code),
		Data:    data,
	}

	c.JSON(httpStatus, r)
}
