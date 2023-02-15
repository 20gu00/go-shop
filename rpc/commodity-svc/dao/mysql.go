package dao

import (
	"commodity-rpc/common/setUp/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	// 全局变量
	DB *gorm.DB
)

func InitMysql(cfg *config.MysqlConfig) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		cfg.UserName,
		cfg.MysqlPassword,
		cfg.MysqlAddr,
		cfg.MysqlPort,
		cfg.DBName,
	)

	// 建立连接
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	//DB.AutoMigrate(&model.Commodity{}, &model.Brand{}, &model.Banner{}, &model.CategoryBrand{}, &model.Category{})
}
