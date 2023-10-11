package router

import (
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

	return router, nil
}
