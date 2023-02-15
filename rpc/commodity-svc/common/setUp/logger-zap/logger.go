package logger

import (
	"commodity-rpc/common/setUp/config"
	"fmt"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// Init 初始化lg
func InitLogger(cfg *config.LogConfig, mode string) (err error) {
	writeSyncer := getLogWriter(cfg.Filename, cfg.MaxSize, cfg.MaxBackup, cfg.MaxAge, cfg.Compress)
	encoder := getEncoder()

	var level = new(zapcore.Level)
	err = level.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return
	}

	var core zapcore.Core
	if mode == "dev" {
		// 进入开发模式，日志输出到终端
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()) //text
		//多个输出
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, level),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	} else if mode == "prod" {
		core = zapcore.NewCore(encoder, writeSyncer, level)
	} else {
		fmt.Println("程序配置的模式不正确")
	}

	logger = zap.New(core, zap.AddCaller())

	//全局变量使用不方便,可以替换zap的全局变量logger
	zap.ReplaceGlobals(logger) //可有可无,在配置logger是最好使用这里的logger配置,外部使用可以使用zap的logger
	zap.L().Info("初始化logger完成")
	//sugarLogger := zap.L().Sugar()
	return nil
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig) //json方便日志收集处理

	//if config.Conf.Mode == "dev" {
	//	return zapcore.NewConsoleEncoder(encoderConfig)
	//} else if config.Conf.Mode == "prod" {
	//	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	//} else {
	//	fmt.Println("程序配置的模式不正确")
	//}

	//zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())  text风格
	//zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()) //json console mapObject

}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int, compress bool) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
		Compress:   compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// GinLogger 接收gin框架默认的日志
//func GinLogger() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		start := time.Now()
//		path := c.Request.URL.Path
//		query := c.Request.URL.RawQuery
//
//		c.Next()
//
//		cost := time.Since(start)
//		//logger
//		logger.Info(path,
//			zap.Int("status", c.Writer.Status()),
//			zap.String("method", c.Request.Method),
//			zap.String("path", path),
//			zap.String("query", query),
//			zap.String("ip", c.ClientIP()),
//			zap.String("user-agent", c.Request.UserAgent()),
//			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
//			zap.Duration("cost", cost),
//		)
//	}
//}
//
//// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
//func GinRecovery(stack bool) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		defer func() {
//			if err := recover(); err != nil {
//				// Check for a broken connection, as it is not really a
//				// condition that warrants a panic stack trace.
//				var brokenPipe bool
//				if ne, ok := err.(*net.OpError); ok {
//					if se, ok := ne.Err.(*os.SyscallError); ok {
//						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
//							brokenPipe = true
//						}
//					}
//				}
//
//				httpRequest, _ := httputil.DumpRequest(c.Request, false)
//				if brokenPipe {
//					logger.Error(c.Request.URL.Path,
//						zap.Any("error", err),
//						zap.String("request", string(httpRequest)),
//					)
//					// If the connection is dead, we can't write a status to it.
//					c.Error(err.(error)) // nolint: errcheck
//					c.Abort()
//					return
//				}
//
//				if stack {
//					logger.Error("[Recovery from panic]",
//						zap.Any("error", err),
//						zap.String("request", string(httpRequest)),
//						zap.String("stack", string(debug.Stack())),
//					)
//				} else {
//					logger.Error("[Recovery from panic]",
//						zap.Any("error", err),
//						zap.String("request", string(httpRequest)),
//					)
//				}
//				c.AbortWithStatus(http.StatusInternalServerError)
//			}
//		}()
//		c.Next()
//	}
//}
