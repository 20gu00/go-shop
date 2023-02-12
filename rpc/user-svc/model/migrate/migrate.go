/*
	用于生成表,而不是放在程序代码中生成,逻辑更加清晰
*/
package main

import (
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"log"
	"os"
	"time"
	"user-rpc/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// create database go_shop default character set utf8mb4 collate utf8mb4_unicode_ci;
	// Dl123456
	//  docker run -p 13306:3306 --name my-mysql -v $PWD/conf:/etc/mysql -v $PWD/data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=Dl123456 -d mysql:5.7

	dsn := "root:Dl123456@tcp(127.0.0.1:13306)/go_shop?charset=utf8&parseTime=True&loc=Local"
	gormLogger := logger.New(
		// io.Writer
		// interface{} any 可以忽略
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			// 慢日志阈值
			SlowThreshold: time.Second,
			// 日志级别
			LogLevel: logger.Info,
			// 是否彩色打印
			Colorful: true,
		},
	)

	// 建立连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		//NamingStrategy: schema.NamingStrategy{
		//	SingularTable: true, // 不用加s后缀
		//},
	})
	if err != nil {
		panic(err)
	}

	// 默认users
	db.AutoMigrate(model.User{})

	options := &password.Options{16, 100, 30, sha512.New}
	salt, encodedPwd := password.Encode("admin12345", options)
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	// 批量注入写数据
	for i := 0; i < 10; i++ {
		user := model.User{
			NickName: fmt.Sprintf("用户%d", i),
			Mobile:   fmt.Sprintf("12345678%d", i),
			Gender:   "Male",
			Password: newPassword,
		}
		db.Save(&user)
	}
}
