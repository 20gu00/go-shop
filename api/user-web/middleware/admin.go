package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"user-web/common/jwt"
)

// 验证用户是否是管理员,就是用于对部分接口添加管理员才能使用的权限

func AdminMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		claim, _ := c.Get("claim")
		user := claim.(*jwt.MyClaims)
		// 2 admin
		if user.AuthId != 2 {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "没有权限",
			})
			c.Abort()
			return
		}
		c.Next()
	}

}
