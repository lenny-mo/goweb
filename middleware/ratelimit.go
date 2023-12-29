package middleware

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	channel chan struct{}
	once    sync.Once // 限制开启漏水口
}

// 创建两个 channel
var (
	// 每个路由组都有自己的单独限流器, 路由组name作为key, channel 作为value
	ratelimitMap = sync.Map{}
)

// RateLimit
//
// ctx 用来控制漏水口的结束
//
// chanName 指定队列的名称，capacity 指定队列长度，1秒内最多可以承受的最大并发数
func RateLimit(ctx context.Context, chanName string, capacity int64) func(c *gin.Context) {
	// 1. 创建一个限流队列，并且开启一个对应的出水口
	// 如果这个键存在，就返回已存在的值；如果这个键不存在，就存储一个新的值（新创建的带有指定容量的通道）
	rateLimiterIface, loaded := ratelimitMap.LoadOrStore(chanName, &rateLimiter{
		channel: make(chan struct{}, capacity),
	})
	limiter := rateLimiterIface.(*rateLimiter) // 指针类型可以直接修改结构体内部的值，而不需要进行值的拷贝
	if !loaded {
		perRequest := time.Second * 10 / time.Duration(capacity) // 每个请求的间隔
		fmt.Println("每个请求的间隔: ", perRequest)
		limiter.once.Do(func() {
			fmt.Println("开启排水功能：", chanName)
			go getOutOfMyBucket(ctx, chanName, limiter.channel, perRequest) // 2. 开启一个出水口
		})
	}

	// 3. 放入请求
	return func(c *gin.Context) {
		// 根据名字获取到对应的通道
		value, _ := ratelimitMap.Load(chanName)
		ch := value.(*rateLimiter).channel

		select {
		case ch <- struct{}{}:
			// 如果能在限时时间内写入则 继续执行
			c.Next()
		case <-time.After(200 * time.Millisecond):
			// 如果时间到了，则直接返回限流信息
			c.AbortWithStatusJSON(429, gin.H{"message": "Too Many Requests"})
		}
	}
}

// getOutOfMyBucket 定时排出队列中的水
func getOutOfMyBucket(ctx context.Context, name string, c <-chan struct{}, perRequest time.Duration) {
	// 开启一个goroutine 执行漏水任务
	for {
		select {
		case <-ctx.Done():
			// 父coroutine退出
			fmt.Println("漏水任务安全退出: ", name)
			return
		default:
			// 定时排水任务
			<-time.After(perRequest)
			<-c
			fmt.Println("排水成功", "队列容量：", len(c))
		}
	}

}
