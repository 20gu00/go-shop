package initdo

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
	"user-web/common"
	"user-web/common/global"
	"user-web/common/nacos"
	"user-web/common/validators"
	"user-web/dao/redis"
	"user-web/pb"

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

	// 初始化全局的user rpc client
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", config.Conf.ConsulConfig.Host, config.Conf.ConsulConfig.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		zap.L().Error("[ GetUserList ]创建consul的client失败")
		return
	}

	userRpcHost := ""
	userRpcPort := 0
	//data,err:=client.Agent().ServicesWithFilter(`Service == "user-rpc"`)
	// 这个格式很重要  或者转义比如\"
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service =="%s"`, config.Conf.UserRpcConfig.Name))
	if err != nil {
		zap.L().Error("[ GetUserList ]从consul总过滤服务失败")
		return
	}
	for _, v := range data {
		userRpcHost = v.Address
		userRpcPort = v.Port
		// 获取这个service任意一个负载即可
		break
	}

	if userRpcHost == "" {
		zap.L().Error("[ GetUserList ]获取rpc服务负载实例失败")
		return
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", userRpcHost, userRpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Error("连接grpc server失败", zap.Error(err))
	}
	global.GlobalUserClient = pb.NewUserClient(conn)

	// 如果user rpc服务下线了,改ip或者port,那么这个客户端也要需改,这个可以在负载均衡中做
	// 好处就是建立好一个客户端,可以直接使用,后续不用在去建立tcp三次握手连接(http2也是基于tcp)
	// 问题,一个连接多个goroutine来使用会有性能问题,于是可以考虑做连接池
	// grpc-connection-pool或者grpc-go-pool
	// 可以自己根据这两个项目做个连接池,或者直接使用consul做负载均衡

	// 从nacos读取配置
	_ = nacos.GetConfigFromNacos()

	<-ch
	return
}
