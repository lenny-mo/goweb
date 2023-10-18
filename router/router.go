package router

import (
	"go_web_app/controllers"
	"go_web_app/logger"

	"github.com/gin-gonic/gin"
)

func Init() (*gin.Engine, error) {
	router := gin.New()
	router.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 测试路由
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})

	// 注册业务路由
	router.POST("/signup", controllers.SignUpHandler)

	return router, nil
}
