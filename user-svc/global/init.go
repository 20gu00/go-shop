package global

import "go-shop/user-svc/dao"

func init() {
	dao.InitMysql()
}
