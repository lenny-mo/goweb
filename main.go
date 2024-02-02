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
	"sync"
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
	// 在程序退出的时候调用sync方法把缓存区日志罗盘
	defer zap.L().Sync()

	// 从终端接收配置文件路径
	var configFile string
	flag.StringVar(&configFile, "c", "./conf/config.yaml", "配置文件的路径")
	flag.Parse()

	// 1. 初始化配置文件
	if err := settings.Init(configFile); err != nil {
		fmt.Println("Init settings failed, err: ", err)
		panic(err)
	}
	fmt.Println(settings.Config.ToString())

	// --------------------- 后续的配置都是在拿到配置文件结构体后进行， 这里使用了BDD：不要隐式引用外部依赖
	//
	// 2. log文件初始化
	if err := logger.Init(settings.Config.LogConfig, settings.Config.Mode); err != nil {
		fmt.Println("Init logger failed, err: ", err)
		panic(err)
	}

	// 3. mysql数据库初始化
	fmt.Println("mysql数据库初始化")
	if err := mysql.Init(settings.Config.MySQLConfig); err != nil {
		fmt.Println("Init mysql failed, err: ", err)
		// 打印调用栈信息
		fmt.Println("Stacktrace from panic: \n" + string(debug.Stack()))
		panic(err)
	}
	defer mysql.Close()

	// 4. redis数据库初始化
	fmt.Println("redis初始化")
	if err := redis.Init(settings.Config.RedisConfig); err != nil {
		fmt.Println("Init redis failed, err: ", err)
		panic(err)
	}
	defer redis.Close()

	// 5. 注册路由
	ctx, cancel := context.WithCancel(context.Background()) // 主函数退出的时候，停止路由中的goroutine
	defer cancel()
	router, err := router.Init(ctx)
	if err != nil {
		fmt.Println("Init router failed, err: ", err)
		zap.L().Error(fmt.Sprintf("Init router failed, err: %v", err))
		panic(err)
	}

	// 6. 注册雪花算法
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	if err := snowflake.Init(formattedTime, settings.Config.MachineId); err != nil {
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
			zap.L().Error(err.Error())
		}
	}()

	// 实现优雅关机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 使用waitgroup 等待优雅关闭结束
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		<-quit
		zap.L().Info("Shutdown Server ...")

		// 关闭服务
		ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelTimeout()

		if err := server.Shutdown(ctxTimeout); err != nil {
			zap.L().Error("Server shutdown failed, err: ", zap.Error(err))
		}
		zap.L().Info("Server exit")
	}(&wg)

	wg.Wait()
}
