package global

import (
	"fmt"
	"go.uber.org/zap"
	"time"
	"user-rpc/common/setUp/config"
	"user-rpc/common/setUp/logger-zap"
	"user-rpc/dao"
	"user-rpc/model"
)

func Init() {
	//初始化logger
	if err := logger.InitLogger(config.Conf.LogConfig, config.Conf.Mode); err != nil {
		fmt.Printf("初始化logger失败, err:%v\n", err)
		panic(err)
	}

	//Timer定时写入
	// 主goroutine退出会影响,看什么方式运行Init,main的goroutine是主线程
	go func() {
		<-time.NewTimer(1 * time.Hour).C
		zap.L().Sync()
	}()
	//defer zap.L().Sync() //写入磁盘
	dao.InitMysql(config.Conf.MysqlConfig)
	dao.DB.AutoMigrate(model.User{})
	//<-ch
	//return
}
