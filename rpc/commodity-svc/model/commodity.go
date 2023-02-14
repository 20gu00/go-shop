package model

import (
	"gorm.io/gorm"
	"time"
)

// proto没有int 由int32 int64
type Base struct {
	ID        int32     `gorm:"primarykey;type:int"` // 如果由外键使用这个主键做关联,那么类型要一致  bigint
	CreatedAt time.Time `gorm:"column:add_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt
	IsDelete  bool
}

// 商品表
type Commodity struct {
}

// 商品分类表   主流电商系统都是三级分类 家用电器->电视->高清电视,每一级分类都有自己的信息,如果每一集定义一张表也就是个结构体,那么后续改动可能就要修改数据库
// 每个分类级别都有自己的信息,一张张表
type Category struct {
	Base
	// 30个字符   非空(不加not null默认可以为空)
	Name string `gorm:"type:varchar(30);not null"`
	// 能否设置在tab栏目(就是上行的,有全部分类点开,也有写已有的分类显示在tab行)
	IsTop bool  `gorm:"default:false;not null"`
	Level int32 `gorm:"type:int;default:1;not null"` // 1 2 3

	// 父级定义 外键
	ParCategoryId int32
	// 不可以循环嵌套
	// 如果是自身,要使用指针
	ParCategory *Category
}

// 品牌表
type Brand struct {
	Base
	Name string `gorm:"type:varchar(30);not null"`
	Logo string `gorm:"type:varchar(300);default:'';not null"` // 不填就默认空字符串

	// 品牌要跟商品分类建立联系,不然选择一个商品分类时要指定品牌的话,不建立联系就需要从全部的品牌中选择相应品牌
	// 品牌和商品分类之间的多多对关系,比如一个商品分类手机,品牌有华为小米,每个品牌有有多个商品分类,所以是品牌表和分类表的多对多关系的中间表

}

// gorm可以自动根据tag生成多对多的中间表
// 自行手动建立
type CategoryBrand struct {
	Base
	// 两个外键来关联两张表
	// 建立为宜索引和一个普通索引,要给这两个外键建立索引
	CategoryId int32 `gorm:"type:int;index:idx_category_brand;unique"`
	Category   Category

	BrandId int32 `gorm:"type:int;index:idx_category_brand;unique"`
	Brand   Brand
}

// 默认生成表明加s,如果是驼峰那就是加下划线
// 自定义表名
func (CategoryBrand) TableName() string {
	return "categorybrand"
}

// 轮播图,可以购买更改
type Banner struct {
	Base
	Picture string `gorm:"type:varchar(300);not null"`
	Url     string `gorm:"type:varchar(300);not null"`
	Index   int32  `gorm:"type:int;default:1;not null"`
}
