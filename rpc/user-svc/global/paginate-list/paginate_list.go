package paginate_list

import (
	"gorm.io/gorm"
)

func Paginate(pNum, pSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pNum == 0 {
			pNum = 1
		}

		// 可以不设置
		switch {
		// 每页数据多余一百,那么就每页大小一百,做个限制每页大小最多给一百
		case pSize > 100:
			pSize = 100
		case pSize <= 0:
			pSize = 10
		}
		
		offset := (pNum - 1) * pSize
		// limit a列出多少   limit(a,b)从哪里列出,列出多少
		return db.Offset(offset).Limit(pSize)
	}
}
