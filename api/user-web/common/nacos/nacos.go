package nacos

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go.uber.org/zap"
	"user-web/common/setUp/config"
)

// 可以配置多个ServerConfig，客户端会对这些服务端做轮询请求
// 至少一个server config
func GetConfigFromNacos() error {
	c := config.Conf
	serverConfig := []constant.ServerConfig{
		{
			//ContextPath: //string // Nacos的ContextPath，默认/nacos，在2.0中不需要设置
			IpAddr: c.NacosConfig.Host,         //"192.168.23.146", //string // Nacos的服务地址
			Port:   uint64(c.NacosConfig.Port), //8848,             //uint64 // Nacos的服务端口
			//Scheme, string     // Nacos的服务地址前缀，默认http，在2.0中不需要设置
			//GrpcPort    uint64 // Nacos的 grpc 服务端口, 默认为 服务端口+1000, 不是必填
		},
	}

	clientConfig := constant.ClientConfig{
		// nacos相关的namespace的id
		NamespaceId:         c.NacosConfig.Namespace, //"e525eafa-f7d7-4029-83d9-008937f9d468", // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		// 最好手动提前准备好这两个目录
		LogDir:   "tmp/nacos/log", // 不用/开头,放在当前项目目录下
		CacheDir: "tmp/nacos/cache",
		LogLevel: "debug",
	}

	// 创建动态配置客户端
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfig,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	//获取配置
	//读取处配置文件,是字符串,比如yaml字符串(json转struct简单,如何实现yaml转struct)
	//方法就是将我们想用的yaml文件去在线yaml转json网址转换成json然后写入nacos,那么从nacos中获取的就是json了,再转换成yaml
	content, err := configClient.GetConfig(vo.ConfigParam{
		// 配置集合
		DataId: c.NacosConfig.DataId, //"user-api"服务自身的配置信息  "user-rpc",
		Group:  c.NacosConfig.Group,  //"dev"
	})

	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	//fmt.Println(content)
	appConfig := new(config.AppConfig)
	// 注意要将这个json字符串转换成这个appConfig,那么这个struct要设置json的tag,不设置可能转换后为空
	json.Unmarshal([]byte(content), &appConfig)
	fmt.Println(appConfig)

	//监听配置文件变化
	err := configClient.ListenConfig(vo.ConfigParam{
		DataId: c.NacosConfig.DataId, //"user-rpc",
		Group:  c.NacosConfig.Group,  //"dev",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("配置文件发生变化")
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	})
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	//而且配置文件是会缓存进cache目录,data_id等,那么如果你修改了上面的nacos的server config的ip或者port等导致不能正确访问nacos,还是会从缓存中拿到配置文件
	//同样如果nacos变更了ip和port短时间不会影响配置文件的获取
	return nil
}
