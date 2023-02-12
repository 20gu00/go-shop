package redis

import (
	"time"
	"user-web/common/setUp/config"
)

func InSmsCode(key, value string) {
	// 两分钟
	rdb.Set(key, value, time.Duration(config.Conf.Exp)*time.Second)
}
