package router

import (
	"context"
	"go_web_app/controllers"
	_ "go_web_app/docs"
	"go_web_app/logger"
	"go_web_app/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// 传入ctx 控制限流器的漏水口速度
func Init(ctx context.Context) (*gin.Engine, error) {
	router := gin.New()
	pprof.Register(router) // 开启性能监控
	router.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册接口文档路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 注册业务路由组, user
	userGroup := router.Group("/user")
	userGroup.Use(middleware.RateLimit(ctx, "user", 100)) // 这个组下面的接口都需要经过限流
	{
		userGroup.POST("/signup", controllers.SignUpHandler)
		userGroup.POST("/login", controllers.LoginHandler)
	}

	// 给前端暴露获取community 分类的接口
	communityGroup := router.Group("/community")
	communityGroup.Use(middleware.RateLimit(ctx, "community", 2))
	communityGroup.Use(middleware.JWT()) // 这个组下面的接口都需要经过jwt验证
	{
		// 业务路由，在指定的社区下创建post
		communityGroup.POST("/:id/createpost", controllers.CreatePostHandler)
		// 业务路由, 通过帖子的id访问帖子详情
		communityGroup.GET("/post/:post_id", controllers.GetPostDetailHandler)
		// 给帖子投票
		communityGroup.POST("/post/vote", controllers.PostVoteHandler)
		// 根据社区id, 时间or 分数来对post list 进行排序
		communityGroup.GET("/:id/sortedpost", controllers.CommunitySortedPostHandler)
	}

	return router, nil
}
