package middleware

import (
	"net/http"
	"strings"
	"user-web/common/jwt"

	"github.com/gin-gonic/gin"
)

// JWT认证中间件
// 先登录获取token才能访问其他接口
func JWTMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		loginReqUrl := "/api/v1/login"
		if len(c.Request.URL.String()) >= len(loginReqUrl) && c.Request.URL.String()[0:13] == loginReqUrl {
			c.Next()
		} else {
			// 携带Token有三种方式 1.放在请求头(header中自定义key value  token:xxx 2.放在请求体 3.放在URI
			// (authorization bear token Token)放在Header的Authorization中，并使用Bearer开头 Authorization: Bearer xxx  / X-TOKEN: xxx
			authHeader := c.Request.Header.Get("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"msg": "请登录,未携带token",
				})
				c.Abort() //ctx不在向下传递(request response)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if !(len(parts) == 2 && parts[0] == "Bearer") {
				c.JSON(http.StatusUnauthorized, gin.H{
					"msg": "传递的token不正确",
				})
				c.Abort()
				return
			}

			//验证token是否有效
			mc, err := jwt.ParseToken(parts[1])
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"msg": "无效的token",
				})
				c.Abort()
				return
			}

			//if err := jwt.OneTokenIng(string(mc.UserID), parts[1]); err != nil {
			//	if err != nil {
			//		common.RespErr(c, common.CodeTwoDevice)
			//		c.Abort()
			//		return
			//	}
			//}

			// 将当前请求的userID信息保存到请求的上下文c上
			// 如果采用session,往往会将用户信息sessionInfo
			c.Set("user_id", mc.ID) // token claim的id
			c.Set("claim", mc)

			c.Next()
		}
	}
}
