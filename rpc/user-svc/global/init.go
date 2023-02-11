package global

import (
	"go-shop/dao"
	"go-shop/model"
)

func Init() {
	dao.InitMysql()
	dao.DB.AutoMigrate(model.User{})
}
