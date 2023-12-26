# golang实现一个发帖论坛后端

## 架构图




## 技术亮点

### channel实现限流中间件

#### 1. 数据结构

 a. `ratelimitMap` (sync.Map)
- 作用：存储每个路由组的限流器
- 结构：使用 `sync.Map`，以路由组名作为键，对应的 Channel 作为值
- 功能：确保每个路由组都有自己的独立限流器，隔离不同路由组的请求

b. `leakDropletMap` (sync.Map)
- 作用：存储每个限流器对应的漏水口
- 结构：使用 `sync.Map`，以路由组名作为键，对应的 `sync.Once` 对象作为值
- 功能：保证每个限流器都有独立的漏水口，并且只初始化一次

```go
var (
	// 每个路由组都有自己的单独限流器, 路由组name作为key, channel 作为value
	ratelimitMap = sync.Map{}
	// 存储每个限流器对应的漏水口, value 是一个sync.Once对象
	leakDropletMap = sync.Map{}
)
```


#### 2. RateLimit 主函数
##### a. RateLimit(chanName string, capacity int64) func(c *gin.Context)
- 参数：
  - `chanName`: 路由组名，用于标识不同的限流器
  - `capacity`: 限流器容量，表示能够同时处理的请求数量
- 作用：返回一个中间件函数，用于限制请求流量
- 过程：
  1. 创建一个 Channel (`createChanByName`)
  2. 开启一个漏水口 (`getOutOfMyBucket`)
  3. 执行请求处理：如果队列没有满则直接写入队列并且使用c.Next执行后续任务；如果队列满了则等待一段时间，还没有获得写入资格则丢弃此次请求


#### 3. 函数详解
##### a. createChanByName(name string, capacity int64)
- 参数：
  - `name`: 路由组名
  - `capacity`: 限流器容量
- 作用：根据名称创建对应的限流队列
- 过程：确保每个路由组有自己的限流队列，若不存在则创建

##### b. getOutOfMyBucket(name string, perRequest time.Duration)
- 参数：
  - `name`: 路由组名
  - `perRequest`: 每个请求的耗时
- 作用：根据名称创建一个单例函数，设置流水速度，定时排出队列中的水
- 过程：
  1. 确保漏水口只创建一次
  2. 获取对应队列，开启 goroutine 执行漏水任务，以控制请求的流出速度

#### 4. 工作原理
- `RateLimit` 函数返回一个中间件，使用选择语句（`select`）控制请求的写入和超时返回。
- 每个路由组通过其名称对应一个独立的 Channel，控制并发请求的数量。
- 漏水口根据每个请求的耗时定时排出队列中的请求，确保在设定的时间内控制请求的流量。