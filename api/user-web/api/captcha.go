package api

/*
	直接在api层做图片验证码
*/

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"net/http"
)

// 验证码的存储
var (
	captchaStore = base64Captcha.DefaultMemStore
)

// 获取验证码
func GetCaptcha(ctx *gin.Context) {
	// 数字 中文 string 动态
	// 长宽像素等的设置
	d := base64Captcha.NewDriverDigit(100, 250, 10, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(d, captchaStore)
	// 生成的验证码的id 编码
	id, b, err := captcha.Generate()
	if err != nil {
		zap.S().Errorf("生成验证码失败:", err.Error())
		ctx.JSON(http.StatusInsufficientStorage, gin.H{
			"msg": "生成验证码错误",
		})
		return
	}
	// 这些信息填充进前段即可显示
	ctx.JSON(http.StatusOK, gin.H{
		// 图片验证码的id
		"captchaId": id,
		// 图片路径,实际上是base64编码的信息,可以在线解码试试就能看到图片验证码
		"picturePath": b,
	})
}
