package main

import (
	"fmt"
	"go_web_app/dao/mysql"
	"go_web_app/dao/redis"
	"go_web_app/logger"
	"go_web_app/settings"
)

func main() {
	// 1. 初始化配置文件
	if err := settings.Init(); err != nil {
		fmt.Println("Init settings failed, err: ", err)
		panic(err)
	}

	// 2. log文件初始化
	if err := logger.Init(); err != nil {
		fmt.Println("Init logger failed, err: ", err)
		panic(err)
	}

	// 3. mysql数据库初始化
	if err := mysql.Init(); err != nil {
		fmt.Println("Init mysql failed, err: ", err)
		panic(err)
	}

	// 4. redis数据库初始化
	if err := redis.Init(); err != nil {
		fmt.Println("Init redis failed, err: ", err)
		panic(err)
	}

	// 5. 注册路由

	fmt.Println("Config init success")

}
