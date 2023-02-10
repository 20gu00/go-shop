package model

import (
	"gorm.io/gorm"
	"time"
)

// 公用字段 不使用gorm的默认model
type Base struct {
	ID int32 `gorm:"primarykey"`
	// 这三个时间在数据库中生成的表都是datetime
	CreatedAt time.Time      `gorm:"column:add_time"` // add_time 数据库字段列名称
	UpdatedAt time.Time      `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt // 默认列名  gorm.Model

	// 不用gorm操作,仅在代码中使用,标记是否删除
	IsDelete bool
}

// 用户表
// 默认生成列名 _
type User struct {
	Base
	// 可以根据mobile查找,创建个索引针对表数据多的情况快速索引 11位
	// 两个索引 normal索引idx_mobile和unique索引
	Mobile string `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	// 加盐加密  对称加密 非对称加密(公钥私钥)  md5摘要算法(不可反解)(即使是开发者也不知道用户实际密码,只能输入密码加密后比对)
	// 不加盐可以暴力破击,就是通过一张彩虹表将常见的密码和对应的md5值记录下来
	Password string `gorm:"type:varchar(250);not null"`
	// 昵称可以为空
	NickName string `gorm:"type:varchar(20)"`
	// 指针,不容易gorm出错
	Birthday *time.Time `gorm:"type:datetime"`
	// type后面加comment
	Gender string `gorm:"column:gender;default:male;type:varchar(6) comment 'male表示男性,female表示女性'"`
	// 简单权限
	Role int `grom:"column:role;default:1;type:int comment '1表示普通用户,2表示管理员'"`
}
