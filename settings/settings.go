package settings

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/fsnotify/fsnotify" // 监控文件

	"github.com/spf13/viper"
)

// Init 把项目的配置文件读取到 viper 中
func Init(configFileName string) error {

	// 1. 通过viper 读取配置文件信息
	// 通过命令行传入的配置文件名读取配置文件信息
	viper.SetConfigFile(configFileName)

	// 2. 两个命令搭配使用
	//viper.SetConfigName("config") // 设置配置文件的名称, 路径在项目根目录下的config.yaml，不包括后缀
	//viper.AddConfigPath(".")      // 添加配置文件的查找路径, 项目的根目录，可以添加多个路径

	err := viper.ReadInConfig() // 读取配置文件内容
	if err != nil {
		// 如果读取配置文件出错，则输出错误信息
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件没有找到，可以忽略错误
			return fmt.Errorf("配置文件没有找到: %s \n", err)
		}
		return fmt.Errorf("Fatal error config file: %s \n", err)
	}

	// 2. 通过读取的信息反序列化到结构体
	if err := viper.Unmarshal(Config); err != nil {
		fmt.Println("Init config file fail!")
	}

	// 3. 监控配置文件的变化
	viper.WatchConfig() // 监控配置文件变化, 如果文件发生变化, 则热加载配置文件
	// fsnotify.Event 可以包含文件被修改的信息
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name) // e.Name 返回被修改的文件名
		zap.L().Info(fmt.Sprintf("Config file %s has been changed:", e.Name))
		fmt.Println("Operation type:", e.Op) // e.Op 返回对文件进行的操作
		zap.L().Info(fmt.Sprintf("Operation type %s has been done on this  file", e.Op.String()))
		if err := viper.Unmarshal(Config); err != nil {
			fmt.Println("Update config file")
			zap.L().Info(fmt.Sprintf("Update config file %s", e.Name))
		}
	})

	return nil
}
