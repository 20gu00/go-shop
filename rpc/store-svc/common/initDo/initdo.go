package initDo

import (
	"fmt"
	"go.uber.org/zap"
	"store-rpc/common/setUp/config"
	"store-rpc/common/setUp/logger-zap"
	"store-rpc/dao"
	"store-rpc/model"
	"time"
)

func Init() {
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
	dao.DB.AutoMigrate(&model.Inventory{})
	//<-ch
	//return
}
