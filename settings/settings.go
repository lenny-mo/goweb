package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

func Init() error {
	// 设置配置文件的名称（不带文件扩展名）
	viper.SetConfigName("config")

	// 设置配置文件的类型为YAML
	viper.SetConfigType("yaml")

	// 添加配置文件的查找路径, 项目的根目录
	viper.AddConfigPath("../")

	// 读取配置文件内容
	err := viper.ReadInConfig()
	if err != nil {
		// 如果读取配置文件出错，则输出错误信息

		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件没有找到，可以忽略错误
			return fmt.Errorf("配置文件没有找到: %s \n", err)
		}

		return fmt.Errorf("Fatal error config file: %s \n", err)
	}

	viper.WatchConfig() // 监控配置文件变化, 如果文件发生变化, 则热加载配置文件
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	return nil
}
