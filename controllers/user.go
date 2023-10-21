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
		ReturnResponse(c, http.StatusBadRequest, InvalidParamCode)
		return
	}
	zap.L().Info("SignUp with param", zap.Any("param", *param)) // 记录结构体信息

	// 2. 业务处理
	err := logic.Signup(param)
	if err != nil {
		zap.L().Error("logic.Signup() failed", zap.Error(err))
		// 3. 返回响应
		ReturnResponse(c, http.StatusBadRequest, InvalidParamCode)
	} else {
		zap.L().Debug("logic.Signup() success")
		ReturnResponse(c, http.StatusOK, SuccessCode)
	}
}

func LoginHandler(c *gin.Context) {
	// 1. 参数校验
	params := new(models.LoginParam)

	if err := c.ShouldBindJSON(params); err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		ReturnResponse(c, http.StatusBadRequest, InvalidParamCode)
		return
	}

	// 2. 业务处理: 判断用户输入的密码是否和数据库中的一致
	err := logic.Login(params)
	if err != nil {
		zap.L().Error("logic.Login() failed", zap.Error(err))
		ReturnResponse(c, http.StatusBadRequest, InvalidParamCode)
	} else {
		zap.L().Debug("logic.Login() success")
		ReturnResponse(c, http.StatusOK, SuccessCode)
	}
}
