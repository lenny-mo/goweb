package middleware

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// 创建两个 channel
var (
	// 每个路由组都有自己的单独限流器, 路由组name作为key, channel 作为value
	ratelimitMap = sync.Map{}
	// 存储每个限流器对应的漏水口, value 是一个sync.Once对象
	leakDropletMap = sync.Map{}
)

// RateLimit
func RateLimit(chanName string, capacity int64) func(c *gin.Context) {
	// 1. 创建一个channel
	createChanByName(chanName, capacity)
	// 2. 开启一个出水口
	perRequest := time.Second / time.Duration(capacity) // 计算每个request的耗时
	fmt.Println("perRequest: ", perRequest)
	getOutOfMyBucket(chanName, perRequest)
	// 3. 放入请求
	return func(c *gin.Context) {
		// 根据名字获取到对应的通道
		value, _ := ratelimitMap.Load(chanName)
		ch := value.(chan struct{})

		select {
		case ch <- struct{}{}:
			// 如果能在限时时间内写入则 继续执行
			c.Next()
		case <-time.After(500 * time.Millisecond):
			// 如果时间到了，则直接返回限流信息
			c.AbortWithStatusJSON(429, gin.H{"message": "Too Many Requests"})
		}
	}
}

// createChanByName 根据名字创建对应的限流队列
func createChanByName(name string, capacity int64) {
	// 判断不存在才创建
	if _, ok := ratelimitMap.Load(name); !ok {
		ratelimitMap.Store(name, make(chan struct{}, capacity))
	}
}

// getOutOfMyBucket 根据名字创建一个单例函数，根据perRequest 设置流水速度
// 根据名字获取到对应的队列，定时排出队列中的水
func getOutOfMyBucket(name string, perRequest time.Duration) {
	// 1. 如果不存在则创建once对象
	if _, ok := leakDropletMap.Load(name); !ok {
		leakDropletMap.Store(name, &sync.Once{})
	}
	// 2. 断言
	value, _ := leakDropletMap.Load(name)
	value.(*sync.Once).Do(func() {
		// 3. 根据名字获取到队列
		fmt.Println("漏水口开始执行：", name)
		c, _ := ratelimitMap.Load(name)
		queue := c.(chan struct{})
		// 开启一个goroutine 执行漏水任务
		go func() {
			for {
				<-time.After(perRequest) // 出水速度：等待10毫秒
				<-queue
				fmt.Println("排水成功")
			}
		}()
	})
}
