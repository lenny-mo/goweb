package main

import (
	"context"
	"fmt"
	"go_web_app/dao/mysql"
	"go_web_app/dao/redis"
	"go_web_app/logger"
	"go_web_app/router"
	"go_web_app/settings"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {

	// 每当完成一个任务，你就记录一个日志。但是，日志并不是立即被保存，而是暂时存放在内存中。
	// 当内存中的日志量达到一定的量时，再将这些日志批量写入到磁盘中。
	// 无论程序是正常结束还是异常退出，都会确保这些日志被保存到文件或其他存储位置。
	defer zap.L().Sync()

	// 1. 初始化配置文件
	if err := settings.Init(); err != nil {
		fmt.Println("Init settings failed, err: ", err)
		panic(err)
	}

	// 2. log文件初始化
	if err := logger.Init(settings.Config.LogConfig); err != nil {
		fmt.Println("Init logger failed, err: ", err)
		panic(err)
	}

	// 3. mysql数据库初始化
	if err := mysql.Init(settings.Config.MySQLConfig); err != nil {
		fmt.Println("Init mysql failed, err: ", err)
		panic(err)
	}
	defer mysql.Close()

	// 4. redis数据库初始化
	if err := redis.Init(settings.Config.RedisConfig); err != nil {
		fmt.Println("Init redis failed, err: ", err)
		panic(err)
	}
	defer redis.Close()

	// 5. 注册路由
	router, err := router.Init()
	if err != nil {
		fmt.Println("Init router failed, err: ", err)
		panic(err)
	}

	// 设置服务器
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: router,
	}

	// 开启goroutine启动服务
	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println("ListenAndServe failed, err: ", err)
			panic(err)
		}
	}()

	quit := make(chan os.Signal)                         // 创建一个无缓冲的通道
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 接收ctrl+c和kill信号
	<-quit

	// 关闭服务
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		zap.L().Error("Server shutdown failed, err: ", zap.Error(err))
	}

	zap.L().Info("Server exit")

}
