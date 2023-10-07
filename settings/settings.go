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

	// 添加配置文件的查找路径
	viper.AddConfigPath(".")

	// 读取配置文件内容
	err := viper.ReadInConfig()
	if err != nil {
		// 如果读取配置文件出错，则输出错误信息
		return fmt.Errorf("Fatal error config file: %s \n", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	return nil
}
