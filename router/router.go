package router

import (
	"go_web_app/controllers"
	"go_web_app/logger"
	"go_web_app/middleware"

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

	// 注册业务路由组, user
	userGroup := router.Group("/user")
	{
		userGroup.POST("/signup", controllers.SignUpHandler)
		userGroup.POST("/login", controllers.LoginHandler)
		userGroup.GET("/index", middleware.JWT(), controllers.IndexHandler) // 只有登录用户可以访问

	}

	// 给前端暴露获取community 分类的接口
	communityGroup := router.Group("/community")
	communityGroup.Use(middleware.JWT()) // 这个组下面的接口都需要经过jwt验证
	{
		communityGroup.GET("/list", controllers.CommunityListHandler)
		// 动态路由
		communityGroup.GET("/:id", controllers.CommunityDetailHandler)
		// 创建post 业务路由
		communityGroup.POST("/:id/post", controllers.CreatePostHandler)
		// 路径参数, 当获通过帖子的id和community id 访问帖子详情时
		communityGroup.GET("/:id/post/:post_id", controllers.GetPostDetailHandler)
	}

	return router, nil
}
