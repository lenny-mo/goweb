package main

import (
	"context"
	"flag"
	"fmt"
	"go_web_app/dao/mysql"
	"go_web_app/dao/redis"
	"go_web_app/logger"
	"go_web_app/pkg"
	"go_web_app/pkg/snowflake"
	"go_web_app/router"
	"go_web_app/settings"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	_ "go_web_app/docs"
)

// @title go_web_app项目接口文档
// @version latest
// @description go_web_app项目接口文档server端api文档
func main() {
	fmt.Println("使用air进行热加载!")

	// 日志并不是立即被保存，而是暂时存放在内存中。
	// 当内存中的日志量达到一定的量时，再将这些日志批量写入到磁盘中。
	// 无论程序是正常结束还是异常退出，都会确保这些日志被保存到文件或其他存储位置。
	defer zap.L().Sync()

	// 从终端接收配置文件路径
	var configFile string
	flag.StringVar(&configFile, "c", "config.yaml", "配置文件的路径")
	flag.Parse()

	// 1. 初始化配置文件
	if err := settings.Init(configFile); err != nil {
		fmt.Println("Init settings failed, err: ", err)
		panic(err)
	}

	// 2. log文件初始化
	if err := logger.Init(settings.Config.LogConfig, settings.Config.Mode); err != nil {
		fmt.Println("Init logger failed, err: ", err)
		panic(err)
	}

	// 3. mysql数据库初始化
	fmt.Println("mysql数据库初始化")
	for i := 0; i < 10; i++ {
		err := mysql.Init(settings.Config.MySQLConfig)
		if err != nil {
			// 睡眠2秒
			time.Sleep(2 * time.Second)
			continue
		} else {
			break
		}
	}
	if err := mysql.Init(settings.Config.MySQLConfig); err != nil {
		fmt.Println("Init mysql failed, err: ", err)
		// 打印调用栈信息
		fmt.Println("Stacktrace from panic: \n" + string(debug.Stack()))
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

	// 6. 注册雪花算法
	if err := snowflake.Init(settings.Config.StartTime, settings.Config.MachineId); err != nil {
		fmt.Println("Init snowflake failed, err: ", err)
		panic(err)
	}

	pkg.Init()

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
