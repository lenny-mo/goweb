package controllers

import (
	"bytes"
	"go_web_app/dao/mysql"
	"go_web_app/settings"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func init() {
	// 初始化数据库的链接， 否则下面的测试会报错
	dbconf := settings.MySQLConfig{
		Host:     "127.0.0.1",
		Port:     3307,
		User:     "root",
		Password: "123456",
		DbName:   "test",
	}
	mysql.Init(&dbconf)
}

func TestSignUpHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	url := "/user/login"
	r.POST(url, LoginHandler)

	body := `
{
    "username": "test",
    "password": "1234"
}
	`
	// 把str 转成reader
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))
	w := httptest.NewRecorder() // 用于记录服务器对 HTTP 请求的响应，以便于测试和断言。
	r.ServeHTTP(w, req)         // 用于模拟请求，把请求的结果记录到 w 中

	// 判断响应码
	assert.Equal(t, http.StatusOK, w.Code)
}
