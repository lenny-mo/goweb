package controllers

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    ResponseCode `json:"code"`
	Message interface{}  `json:"message"`
	Data    interface{}  `json:"data"`
}

func ReturnResponse(c *gin.Context, httpStatus int, code ResponseCode) {
	r := Response{
		Code:    code,
		Message: code.GetMsg(code),
		Data:    nil,
	}

	c.JSON(httpStatus, r)
}
