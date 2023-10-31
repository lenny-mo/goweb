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
		// 获取社区的所有分类
		communityGroup.GET("/list", controllers.CommunityListHandler)
		// 动态路由，通过社区id获取社区详情
		communityGroup.GET("/:id", controllers.CommunityDetailHandler)
		// 业务路由，在指定的社区下创建post
		communityGroup.POST("/:id/createpost", controllers.CreatePostHandler)
		// 业务路由, 通过帖子的id和community id 访问帖子详情时
		communityGroup.GET("/:id/post/:post_id", controllers.GetPostDetailHandler)
		// 获取帖子列表
		communityGroup.GET("/:id/postlist", controllers.GetPostListHandler)
		// 给帖子投票
		communityGroup.POST("/:id/post/vote", controllers.PostVoteHandler)
	}

	return router, nil
}
