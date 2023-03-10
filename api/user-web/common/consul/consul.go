package consulDo

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" //匿名引用 init
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func ConsulRegister(name string, id string, addr string, port int, tag []string) error {
	// 注册的服务和consul的地址最好都不要写127.0.0.1，建议写成192.168.23.100
	// 同名的也就是name注册会覆盖
	// 其实就是ip(port)不一样但是name一样,就是为一个服务注册多个工作服在实例

	cfg := api.DefaultConfig()
	// 这是consul的地址
	cfg.Address = "127.0.0.1:8500" //"192.168.23.100:8500"
	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	// register逻辑
	// 实际上就是发送http请求

	// 要注册的对象
	registor := new(api.AgentServiceRegistration)
	registor.ID = id
	registor.Name = name
	registor.Tags = tag
	// 注册的服务地址
	registor.Address = addr
	registor.Port = port

	// 健康检查的对象
	IpPort := fmt.Sprintf("%s:%d", addr, port)
	registor.Check = &api.AgentServiceCheck{
		HTTP:                           "http://" + IpPort + "/health", // ping
		Timeout:                        "10s",                          // 健康检查超时的时间
		Interval:                       "5s",                           // 检查间隔
		DeregisterCriticalServiceAfter: "10s",                          // 健康检查失败后移除服务

	}
	if err := client.Agent().ServiceRegister(registor); err != nil {
		return err
	}

	return nil
}

// 服务发现
func GetAllSvc() error {
	cfg := api.DefaultConfig()
	// 这是consul的地址
	cfg.Address = "127.0.0.1:8500" //"192.168.23.100:8500"
	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	svcs, err := client.Agent().Services() // Service()拿到指定的单个服务
	if err != nil {
		return err
	}

	for k, v := range svcs {
		fmt.Println(k, v)
	}
	return nil
}

// 获取某个service
func GetService() error {
	cfg := api.DefaultConfig()
	// 这是consul的地址
	cfg.Address = "127.0.0.1:8500" //"192.168.23.100:8500"
	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	// 通过name过滤
	// 获取服务时返回的数据就有Service,就是通过返回的数据的字段来过滤
	// ID Tag等等的
	data, err := client.Agent().ServicesWithFilter(`Service=="user-web"`) // user-web之前注册的服务 name
	if err != nil {
		return err
	}
	for k, v := range data {
		fmt.Println(k, v)
	}
	return nil
}

// 注销服务
func DeRegister(id string) error {
	cfg := api.DefaultConfig()
	// 这是consul的地址
	cfg.Address = "127.0.0.1:8500" //"192.168.23.100:8500"
	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}

	err := client.Agent().ServiceDeregister(id)
	if err != nil {
		return err
	}
	return nil
}

// api层(gin)从consul拉取服务信息
func PullFromConsul() {
	c, err := grpc.Dial(
		// consul中的service_name
		// wait等待解析多长时间
		// limit有多个服务只要多少个
		// tag=manual过滤作用
		// https://github.com/mbobakov/grpc-consul-resolver

		"consul://192.168.23.146:8500/user-rpc?wait=15s",
		// grpc.WithInsecure(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),

		// 负载均衡算法
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// 创建一个grpc client,操作grpc server,这事是直接操作consul返回的信息
	//userRpcFromClient:=pb.NewUserClient(c)
	//可以使用这个client调用方法了
}
