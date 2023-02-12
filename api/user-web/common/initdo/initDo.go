package initdo

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"time"
	"user-web/common"
	"user-web/common/validators"
	"user-web/dao/redis"

	"go.uber.org/zap"

	"user-web/common/setUp/config"
	"user-web/common/setUp/logger-zap"
)

func InitDO(ch chan int) {
	//初始化logger
	if err := logger.InitLogger(config.Conf.LogConfig, config.Conf.Mode); err != nil {
		fmt.Printf("初始化logger失败, err:%v\n", err)
		panic(err)
	}

	//Timer定时写入
	go func() {
		<-time.NewTimer(1 * time.Hour).C
		zap.L().Sync()
	}()
	defer zap.L().Sync() //写入磁盘

	//初始化mysql连接
	//if err := mysql.InitMysql(config.Conf.MysqlConfig); err != nil {
	//	fmt.Printf("初始化mysql失败, err:%v\n", err)
	//	panic(err)
	//}
	// 外部调用的时候使用,或者chan阻塞这个goroutine
	//defer mysql.DBClose()

	//初始化redis连接
	if err := redis.InitRedis(config.Conf.RedisConfig); err != nil {
		fmt.Printf("初始化redis失败, err:%v\n", err)
		panic(err)
	}
	defer redis.RDBClose()

	//雪花算法生成分布式uid
	//if err := snowflake.InitSnowFlake(config.Conf.StartTime, config.Conf.MachineID); err != nil {
	//	fmt.Printf("雪花算法生成uid失败, err:%v\n", err)
	//	panic(err)
	//}

	//初始化gin内置支持的校验器(validator)的翻译器(en zh)
	if err := common.InitTrans("zh"); err != nil {
		fmt.Printf("初始化validator翻译器失败, err:%v\n", err)
		return
	}

	// 自定义的验证规则(验证器)和注册翻译器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 定义tag类似min等
		v.RegisterValidation("mobile", validators.ValidatorMobile)
		// 翻译器,输出提示信息的定义
		_ = v.RegisterTranslation("mobile", common.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 手机号码格式不正确!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	<-ch
	return
}
