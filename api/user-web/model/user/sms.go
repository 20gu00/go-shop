package user

type SmsInput struct {
	// 发送验证码类型,注册(用户不存在就创建用户,存在就返回已经存在),登录(用户存在就登陆,不存在就注册),找回密码
	Kind   string `form:"kind" json:"kind" binding:"required,oneof=register login"`
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"`
}
