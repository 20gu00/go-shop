package router

import (
	"github.com/gin-gonic/gin"
	"user-web/api"
	"user-web/middleware"
)

func SetupRouter(r *gin.Engine) {
	user := r.Group("/user")

	user.GET("/list", middleware.AdminMiddleware(), api.GetUserList)
	user.POST("/user_pwd_login", api.UserPasswdLogin)

	com := r.Group("/com")
	com.GET("/captcha", api.GetCaptcha)
}

