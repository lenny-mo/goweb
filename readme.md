# golang实现一个发帖论坛后端

## 架构图


# 限流组件设计思路

## 概述
本文档描述了一个用于实现限流的组件。该组件基于 Golang 中的并发原理和数据结构来限制每秒内的请求并发数。使用了 `sync.Map` 和 `chan` 数据结构来存储和控制请求的并发量。

## 设计

本次设计主要借鉴了Uber的漏桶算法。

### 数据结构
#### `rateLimiter` 结构体
- `channel chan struct{}`：用于限制请求并发数的通道。每个路由组拥有一个独立的限流通道。
- `once sync.Once`：用于控制排水任务只执行一次的同步机制。

#### `ratelimitMap`（全局变量）
- `sync.Map` 类型，用于存储每个路由组对应的限流器。路由组名称作为键，`rateLimiter` 结构体指针作为值。

### 函数
#### `RateLimit`
- 参数：`ctx` 控制排水任务结束的上下文，`chanName` 指定队列名称，`capacity` 指定队列长度（每秒内最多可以承受的最大并发数）。
- 逻辑：
  1. 创建一个限流队列并开启对应的出水口。
     - 使用 `LoadOrStore` 从 `ratelimitMap` 中加载或存储通道。
     - 若键存在，则返回已存在的值；若键不存在，则存储一个新的值（创建一个带有指定容量的通道）。
     - 确保排水任务只执行一次，并根据容量开启排水任务。
  2. 放入请求：
     - 根据名称获取对应通道，将请求放入通道。
     - 如果能在限定时间内写入，则继续执行请求；否则返回限流信息。

#### `getOutOfMyBucket`
- 参数：`ctx` 控制排水任务结束的上下文，`name` 队列名称，`c` 限流通道，`perRequest` 每个请求的间隔时间。
- 逻辑：
  - 开启一个 goroutine 执行漏水任务，定时排出队列中的水。
  - 当接收到上下文的结束信号时，安全退出漏水任务。
  - 定时排水任务：
    - 每隔一定时间执行排水操作。
    - 从通道中取出一个元素，表示请求处理完毕。
    - 打印排水成功的信息和当前队列容量。

## 用法示例
下面是一个使用该限流组件的示例代码：

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 创建一个上下文，用于控制漏水口的结束
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

  // 绑定名为route1的限流器，每秒最多10个请求
	r.GET("/limited", RateLimit(ctx, "route1", 10), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Request Allowed"})
	})

	r.Run(":8080")
}
```

该示例展示了如何在路由上使用 `RateLimit` 函数来限制请求的并发数。

## 注意事项
- 限流器是基于每个路由组（路由名称）独立创建的，不同路由组之间的限流互不影响。
- 请根据实际需求调整并发量和排水时间间隔，确保合理的限流效果。
- 在实际生产环境中，需考虑更多并发安全性和性能优化方面的问题，例如并发操作的竞态条件等。