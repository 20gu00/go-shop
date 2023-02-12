package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"time"
	"user-web/common/setUp/config"
)

// 保存验证码
func InSmsCode(key, value string) {
	// 两分钟
	rdb.Set(key, value, time.Duration(config.Conf.Exp)*time.Second)
}

// 获取验证码
func GetSmsCode(key string) (string, error) {
	// 通过这个key拿到的结果
	v, err := rdb.Get(key).Result()
	if err != redis.Nil {
		return "", errors.New("key不存在的,或者" + err.Error())
	}
	return v, nil
}
