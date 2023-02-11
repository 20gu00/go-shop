package jwt

import (
	"errors"
	"time"
	//"user-web/dao/redis"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

//签名用的secret
//浏览器在线生成一个随机的密码
var (
	mySecret = []byte("cjq")
)

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 自定义结构体,添加自定义的字段,token加密解密都是这个结构体
type MyClaims struct {
	ID       uint
	AuthId   uint // role admin 普通用户
	Nickname string
	jwt.StandardClaims
}

// 生成JWT
func GenToken(id, authId uint, nickname string) (string, error) {
	// 创建一个自己声明的数据
	c := MyClaims{
		id,
		authId,
		nickname,
		jwt.StandardClaims{
			// NotBefore 生效时间 time.Now().Unix()
			ExpiresAt: time.Now().Add(
				time.Duration(viper.GetInt("auth.jwt_expire")) * time.Minute).Unix(), // 过期时间
			Issuer: "cjq",
		},
	}
	// 使用指定的签名方法创建签名对象(加密算法,token配置)
	//header payload
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 签名并获得完整的编码后的字符串token
	//signature
	return token.SignedString(mySecret)
	//最终的base64编码
}

// 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token,存放字段
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid { // 校验token是否有效
		return mc, nil
	}
	return nil, errors.New("无效的token")
}

// 基于token实现同一个账户只能登陆一台设备(登录状态)
//func OneTokenIng(userID string, token string) error {
//	v, err := redis.GetTokenKey(userID)
//	if err == nil {
//		zap.L().Error("获取userID对应的token失败")
//		return err
//	}
//
//	if v != "" {
//		if v == token {
//			return errors.New("同一个账户一时间内只能登陆一台设备")
//		}
//	}
//
//	return nil
//}
