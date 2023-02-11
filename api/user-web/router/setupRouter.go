package router

import (
	"github.com/gin-gonic/gin"
	"user-web/api"
)

func SetupRouter(r *gin.Engine) {
	user := r.Group("/user")

	user.GET("/list", api.GetUserList)
	user.POST("/user_pwd_login", api.UserPasswdLogin)
}