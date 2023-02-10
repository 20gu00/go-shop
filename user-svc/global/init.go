package global

import (
	"go-shop/user-svc/dao"
	"go-shop/user-svc/model"
)

func Init() {
	dao.InitMysql()
	dao.DB.AutoMigrate(model.User{})
}
