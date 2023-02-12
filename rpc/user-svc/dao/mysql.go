package dao

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"user-rpc/common/setUp/config"
)

var (
	// 全局变量
	DB *gorm.DB
)

func InitMysql(cfg *config.MysqlConfig) {
	//dsn := "root:Dl123456@tcp(127.0.0.1:13306)/go_shop?charset=utf8&parseTime=True&loc=Local"
	//gormLogger := logger.New(
	//	// io.Writer
	//	// interface{} any 可以忽略
	//	log.New(os.Stdout, "\r\n", log.LstdFlags),
	//	logger.Config{
	//		// 慢日志阈值
	//		SlowThreshold: time.Second,
	//		// 日志级别
	//		LogLevel: logger.Info,
	//		// 是否彩色打印
	//		Colorful: true,
	//	},
	//)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		cfg.UserName,
		cfg.MysqlPassword,
		cfg.MysqlAddr,
		cfg.MysqlPort,
		cfg.DBName,
	)

	// 建立连接
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		//Logger: gormLogger,
		// 最好不要
		//NamingStrategy: schema.NamingStrategy{
		//	SingularTable: true, // 不用加s后缀
		//},
	})
	if err != nil {
		panic(err)
	}

	// 默认users
	//DB.AutoMigrate(model.User{})
}
