package controllers

import (
	"go_web_app/logic"
	"go_web_app/models"
	"net/http"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"

	_ "github.com/go-playground/validator/v10"
)

func SignUpHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	param := new(models.SignupParam) // 定义一个结构体变量,内部是默认值
	if err := c.ShouldBindJSON(param); err != nil {
		// 参数有误
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "请求参数有误",
		})
		return
	}
	zap.L().Info("SignUp with param", zap.Any("param", *param)) // 记录结构体信息

	// 2. 业务处理
	err := logic.Signup(param)
	if err != nil {
		zap.L().Error("logic.Signup() failed", zap.Error(err))
		// 3. 返回响应
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "bad request",
		})
	} else {
		zap.L().Debug("logic.Signup() success")
		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
	}
}
