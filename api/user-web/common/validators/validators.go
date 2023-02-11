package validators

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidatorMobile(f validator.FieldLevel) bool {
	// 获取mobile
	mobile := f.Field().String()
	// 手机号码正则表达式
	ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35--8]|9[189])\d{8}$`, mobile)
	if !ok {
		return false
	}
	return true
}
