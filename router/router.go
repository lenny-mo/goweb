package router

import (
	"go_web_app/controllers"
	_ "go_web_app/docs"
	"go_web_app/logger"
	"go_web_app/middleware"
	"math/rand"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/pprof"
)

func Init() (*gin.Engine, error) {
	router := gin.New()
	// 注册pprof路由
	pprof.Register(router)

	router.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册接口文档路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 注册测试路由, 每次随机生成一个数组并且对数组进行快速排序
	router.GET("/ping", pingHandler)

	// 注册业务路由组, user
	userGroup := router.Group("/user")
	userGroup.Use(middleware.RateLimit()) // 这个组下面的接口都需要经过限流

	{
		userGroup.POST("/signup", controllers.SignUpHandler)
		userGroup.POST("/login", controllers.LoginHandler)
		userGroup.GET("/index", middleware.JWT(), controllers.IndexHandler) // 只有登录用户可以访问

	}

	// 给前端暴露获取community 分类的接口
	communityGroup := router.Group("/community")
	communityGroup.Use(middleware.JWT()) // 这个组下面的接口都需要经过jwt验证
	communityGroup.Use(middleware.RateLimit())
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
		// 根据时间or 分数来对post list 进行排序，
		communityGroup.GET("/sortedpost", controllers.SortedPostHandler)
		// 根据社区id, 时间or 分数来对post list 进行排序
		communityGroup.GET("/:id/sortedpost", controllers.CommunitySortedPostHandler)
	}

	return router, nil
}

func pingHandler(c *gin.Context) {
	// 1. 随机生成长度为10000的数组
	list := make([]int, 10000)
	for i := range list {
		list[i] = rand.Intn(10000)
	}
	// 2. 对数组进行快速排序
	quickSort(list, 0, len(list)-1)

	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func quickSort(list []int, low, high int) {
	if low >= high {
		return
	}
	p := partition(list, low, high)
	quickSort(list, low, p-1)
	quickSort(list, p+1, high)
}
func partition(list []int, low, high int) int {
	pivot := list[high] // 选取最后一个元素作为pivot
	i := low            // 使用i指针来找到pivot的位置
	for j := low; j < high; j++ {
		if list[j] < pivot {
			list[i], list[j] = list[j], list[i]
			i++
		}
	}
	list[i], list[high] = list[high], list[i] // pivot 左边的所有元素都小于它，而右边的所有元素都大于或等于它
	return i
}
