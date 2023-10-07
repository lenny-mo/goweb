package main

import (
	"fmt"
	"go_web_app/settings"
)

func main() {
	// 1. 初始化配置文件
	if err := settings.Init(); err != nil {
		panic(err)
	}

	fmt.Println("Config init success")
}
