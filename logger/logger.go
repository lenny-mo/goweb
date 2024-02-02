package logger

import (
	"go_web_app/settings"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Init 初始化日志系统
func Init(conf *settings.LogConfig, mode string) (err error) {
	// 创建一个WriteSyncer（日志写入器），用于写入日志文件
	writeSyncer := getLogWriter(
		conf.Filename,
		conf.MaxSize,
		conf.MaxBackups,
		conf.MaxAge,
	)

	// 获取日志的编码器
	encoder := getEncoder()

	// 定义日志级别
	var level = new(zapcore.Level)

	// 从配置中读取日志级别并解析
	err = level.UnmarshalText([]byte(viper.GetString("log.level")))
	if err != nil {
		return
	}

	// 定义核心日志处理器
	var core zapcore.Core

	// 如果是开发模式，设置控制台和文件双输出
	if mode == "dev" {
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, level),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		// 否则只输出到文件
		core = zapcore.NewCore(encoder, writeSyncer, level)
	}

	// 创建新的zap日志实例
	lg := zap.New(core, zap.AddCaller())

	// 替换zap库的全局实例, 可以通过zap.L()调用
	zap.ReplaceGlobals(lg)

	return
}

// getEncoder 获取日志的编码器，用于定义日志的格式
func getEncoder() zapcore.Encoder {
	// 设置日志的编码配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// 返回JSON格式的编码器
	return zapcore.NewJSONEncoder(encoderConfig)
}

// getLogWriter 创建一个日志写入器
func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	// 使用lumberjack作为日志滚动库
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,  // 日志文件名
		MaxSize:    maxSize,   // 日志最大大小（MB）
		MaxBackups: maxBackup, // 最大备份文件数
		MaxAge:     maxAge,    // 日志保存的最大天数
	}
	// 返回一个WriteSyncer，用于写入日志
	return zapcore.AddSync(lumberJackLogger)
}

// GinLogger 为gin框架提供日志记录功能
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()             // 请求开始时间
		path := c.Request.URL.Path      // 请求路径
		query := c.Request.URL.RawQuery // 请求查询参数
		c.Next()                        // 处理请求

		// 计算请求处理时间
		cost := time.Since(start)
		// 记录日志
		zap.L().Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 检查是否是断开的连接
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					// 如果是断开的连接，记录日志并中止请求
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				// 记录panic信息
				if stack {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next() // 继续处理其他中间件或路由
	}
}
