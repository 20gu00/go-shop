package user

// 用户密码登录请求
// 表单传参数
type UserPasswdLoginForm struct {
	Mobile   string `json:"mobile" form:"mobile" binding:"required,mobile"`           // 自定义手机号码验证规则
	Password string `json:"password" form:"password" binding:"required,min=3,max=50"` // 逗号之间不能有空格,最少3个最大50个
}

// 获取用户列表的响应
type UserListRes struct {
	Id int32 `json:"id"`
	// 赋值时time->string
	Birthday string `json:"birthday"`
	NickName string `json:"nickname"`
	Gender   string `json:"gender"`
	Mobile   string `json:"mobile"`
}
