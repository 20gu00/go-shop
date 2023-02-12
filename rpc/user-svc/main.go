package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"strconv"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"user-rpc/common/setUp/config"
	"user-rpc/global"
	"user-rpc/pb"
	"user-rpc/service"
)

func main() {
	// 返回的是指针
	var confFile string
	flag.StringVar(&confFile, "conf", "", "配置文件")

	// 建议不要写0.0.0.0,写本机ip192.168.23.146,不然consul发现有问题
	// 或者后续的grpc检查的ip从配置中心获取,也就是本机ip
	ip := flag.String("ip", "0.0.0.0", "ip")
	port := flag.String("port", "50051", "port")
	address := fmt.Sprintf("%s:%s", *ip, *port)

	flag.Parse()

	//读取配置文件,加载配置文件需要时间如果用goroutine方式去加载最好主goroutine阻塞一会,不然那拿到的配置值为空
	if err := config.ConfRead(confFile); err != nil {
		fmt.Printf("读取配置文件失败, err:%v\n", err)
		panic(err)
	}

	global.Init()
	fmt.Println("11")
	// 创建server
	server := grpc.NewServer()

	// 注册user service
	pb.RegisterUserServer(server, &service.UserServer{})

	// 注册健康检查
	srv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, srv)

	cfg := api.DefaultConfig()
	// 这是consul的地址
	cfg.Address = fmt.Sprintf("%s:%d", config.Conf.ConsuleConfig.Host, config.Conf.ConsuleConfig.Port) //"127.0.0.1:8500" //"192.168.23.100:8500"
	client, err := api.NewClient(cfg)
	if err != nil {

		return
	}
	// 要注册的对象
	registor := new(api.AgentServiceRegistration)
	registor.ID = config.Conf.Name
	registor.Name = config.Conf.Name
	registor.Tags = []string{"user-rpc"}
	// 注册的服务地址,注意要和下面的GRPC保持一致
	registor.Address = "192.168.23.146"
	p, _ := strconv.Atoi(*port)
	registor.Port = p
	// 健康检查的对象
	//IpPort := fmt.Sprintf("%s:%d", addr, port)
	registor.Check = &api.AgentServiceCheck{
		GRPC:                           "192.168.23.146:50051", //"http://" + IpPort + "/health", // ping
		Timeout:                        "10s",                  // 健康检查超时的时间
		Interval:                       "5s",                   // 检查间隔
		DeregisterCriticalServiceAfter: "10s",                  // 健康检查失败后移除服务

	}
	if err := client.Agent().ServiceRegister(registor); err != nil {
		fmt.Println("22", err.Error())
		return
	}

	// 建立连接
	// ip port 服务提供的地址,监听的地址
	listen, err := net.Listen("tcp", address)
	if err != nil {
		panic("监听连接失败")
	}

	// grpc server启动
	fmt.Println("服务监听的地址是", address)
	err = server.Serve(listen)
	if err != nil {
		panic("server启动失败")
	}
}
