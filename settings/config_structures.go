package settings

// -------------------------- config struct --------------------------- //
var (
	Config = new(Conf)
)

type Conf struct {
	*AppConfig   `mapstructure:"app"`   // AppConfig 是应用的配置。
	*LogConfig   `mapstructure:"Log"`   // LogConfig 是应用的日志配置。
	*MySQLConfig `mapstructure:"mysql"` // MySQLConfig 是应用的 MySQL 数据库配置。
	*RedisConfig `mapstructure:"redis"` // RedisConfig 是应用的 Redis 数据库配置。
}

// AppConfig 是应用的配置结构体。
type AppConfig struct {
	Name      string `mapstructure:"name"`      // Name 是应用的名称。
	Mode      string `mapstructure:"mode"`      // Mode 是应用的运行模式，例如：development, production等。
	Version   string `mapstructure:"version"`   // Version 是应用的版本号。
	Port      int    `mapstructure:"port"`      // Port 是应用监听的端口号。
	StartTime string `mapstructure:"starttime"` // StartTime 是雪花算法的起始时间。
	MachineId int64  `mapstructure:"machineid"` // MachineId 是雪花算法的机器ID。
}

// LogConfig 结构体定义了日志配置的各项参数。
type LogConfig struct {
	Level      string `mapstructure:"level"`      // Level 表示日志的级别（例如：info, warning, error 等）。
	Filename   string `mapstructure:"filename"`   // Filename 表示日志文件的名称。
	MaxSize    int    `mapstructure:"maxsize"`    // MaxSize 表示每个日志文件的最大大小（以MB为单位）。
	MaxAge     int    `mapstructure:"maxage"`     // MaxAge 表示日志文件的最大保留天数。
	MaxBackups int    `mapstructure:"maxbackups"` // MaxBackups 表示保留的最大旧日志文件数。
}

// MySQLConfig 结构体定义了连接 MySQL 数据库所需的配置参数。
type MySQLConfig struct {
	Host         string `mapstructure:"host"`     // Host 表示 MySQL 数据库的主机地址。
	User         string `mapstructure:"user"`     // User 表示连接数据库所使用的用户名。
	Password     string `mapstructure:"password"` // Password 表示连接数据库所使用的密码。
	DbName       string `mapstructure:"dbname"`   // DbName 表示要连接的数据库名称。
	Port         int    `mapstructure:"port"`     // Port 表示 MySQL 数据库的端口号。
	MaxOpenConns int    `mapstructure:"maxopen"`  // MaxOpenConns 表示数据库的最大打开连接数。
	MaxIdleConns int    `mapstructure:"maxidle"`  // MaxIdleConns 表示数据库的最大空闲连接数。
}

// RedisConfig 结构体定义了连接 Redis 数据库所需的配置参数。
type RedisConfig struct {
	Host     string `mapstructure:"host"`     // Host 表示 Redis 数据库的主机地址。
	Password string `mapstructure:"password"` // Password 表示连接数据库所使用的密码。
	Port     int    `mapstructure:"port"`     // Port 表示 Redis 数据库的端口号。
	DB       int    `mapstructure:"db"`       // DB 表示要连接的 Redis 数据库的编号。
	PoolSize int    `mapstructure:"poolsize"` // PoolSize 表示连接池的大小。
}
