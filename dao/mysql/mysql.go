package mysql

import (
	"fmt"
	"go_web_app/settings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// 定义一个全局变量 db，它是一个指向 sqlx.DB 结构体的指针，该结构体保存了数据库连接的所有信息。
var db *sqlx.DB

// initDB 函数负责初始化数据库连接，并返回一个 error 值，以指示是否有任何错误发生。
func Init(conf *settings.MySQLConfig) (err error) {
	// dsn 是数据源名称，它包含了数据库连接所需的所有信息
	// 使用viper 读取配置文件
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DbName,
		//viper.GetString("mysql.user"),
		//viper.GetString("mysql.password"),
		//viper.GetString("mysql.host"),
		//viper.GetInt("mysql.port"),
		//viper.GetString("mysql.dbname"),
	)

	// 使用 sqlx.Connect 函数连接到 MySQL 数据库。
	// 如果连接失败，它将返回一个 error，我们可以检查这个 error 来确定是否成功。
	db, err = sqlx.Connect("mysql", dsn)
	// 如果连接过程中发生错误，打印错误信息并返回 error。
	if err != nil {
		zap.L().Error("connect DB failed, err: ", zap.Error(err))
		return
	}
	// 设置数据库的最大打开连接数为 20。
	db.SetMaxOpenConns(conf.MaxOpenConns)
	// 设置数据库的最大空闲连接数为 10。
	db.SetMaxIdleConns(conf.MaxIdleConns)
	// 函数成功完成，返回 nil（无错误）。
	return
}

func Close() {
	db.Close()
}
