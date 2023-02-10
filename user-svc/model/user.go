package model

import (
	"go-shop/user-svc/dao"
	paginate_list "go-shop/user-svc/global/paginate-list"
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
	// 指针,不容易gorm出错   数据库是datetime
	Birthday *time.Time `gorm:"type:datetime"`
	// type后面加comment
	Gender string `gorm:"column:gender;default:male;type:varchar(6) comment 'male表示男性,female表示女性'"`
	// 简单权限
	Role int `grom:"column:role;default:1;type:int comment '1表示普通用户,2表示管理员'"`
}

// users
//func TableName()string{
//	return "users"
//}

type UserDao interface {
	GetUserList() ([]User, int32, error)
	Paginate(pNum, pSize int) []User
	GetUserByMobile(mobile string) (User, *gorm.DB)
	GetUserById(id int32) (User, *gorm.DB)
	CreateUser(user *User) *gorm.DB
	UpdateUser(user *User) *grom.DB
}

func NewUserDao() UserDao {
	return &User{} // 指针最好,跟接收器一致
}

func (u *User) GetUserList() ([]User, int32, error) {
	var userList []User
	result := dao.DB.Find(&userList)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	total := int32(result.RowsAffected)
	return userList, total, nil
}

func (u *User) Paginate(pNum, pSize int) []User {
	var userPaginate []User
	// 使用gorm的scope
	dao.DB.Scopes(paginate_list.Paginate(pNum, pSize)).Find(&userPaginate)
	return userPaginate
}

func (u *User) GetUserByMobile(mobile string) (User, *gorm.DB) {
	var user User
	// 指针
	result := dao.DB.Where(&User{Mobile: mobile}).First(&user)
	return user, result
}

func (u *User) GetUserById(id int32) (User, *gorm.DB) {
	// 考虑清楚是否直接使用接受者指针
	var user User
	// 指针
	result := dao.DB.Where(&User{Base: Base{ID: id}}).First(&user)
	return user, result
}

func (u *User) CreateUser(user *User) *gorm.DB {
	// 指针
	result := dao.DB.Create(user)
	return result
}

func (u *User) UpdateUser(user *User) *gorm.DB {
	// 指针
	result := dao.DB.Save(user)
	return result
}
