package main

import (
	"commodity-rpc/common/initDo"
	"commodity-rpc/common/setUp/config"
	"commodity-rpc/common/tool"
	"commodity-rpc/pb"
	"commodity-rpc/service"
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	var confFile string
	flag.StringVar(&confFile, "conf", "", "配置文件")
	ip := flag.String("ip", "0.0.0.0", "ip")
	port := flag.String("port", "0", "port")

	if *port == "0" {
		// 没有从命令行获取端口
		p, _ := tool.GetFreePort()
		*port = strconv.Itoa(p)
	}

	address := fmt.Sprintf("%s:%s", *ip, *port)

	if err := config.ConfRead(confFile); err != nil {
		fmt.Printf("读取配置文件失败, err:%v\n", err)
		panic(err)
	}
	initDo.Init()

	server := grpc.NewServer()
	pb.RegisterCommodityServer(server, &service.CommodityServer{})
	srv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, srv)

	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", config.Conf.ConsuleConfig.Host, config.Conf.ConsuleConfig.Port) //"127.0.0.1:8500" //"192.168.23.100:8500"
	client, err := api.NewClient(cfg)
	if err != nil {

		return
	}

	registor := new(api.AgentServiceRegistration)
	// id不同,一个服务多个实例
	id := fmt.Sprintf("%s", uuid.NewV4())
	registor.ID = id //config.Conf.Name
	registor.Name = config.Conf.Name
	registor.Tags = []string{"commodity-rpc"}
	// 注册的服务地址,注意要和下面的GRPC保持一致
	// 可以通过环境变量获取
	registor.Address = "192.168.23.146"
	p, _ := strconv.Atoi(*port)
	registor.Port = p
	// 健康检查的对象
	//IpPort := fmt.Sprintf("%s:%d", addr, port)
	registor.Check = &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("192.168.23.146:%d", *port), //"http://" + IpPort + "/health", // ping
		Timeout:                        "10s",                                   // 健康检查超时的时间
		Interval:                       "5s",                                    // 检查间隔
		DeregisterCriticalServiceAfter: "10s",                                   // 健康检查失败后移除服务

	}
	if err := client.Agent().ServiceRegister(registor); err != nil {
		fmt.Println(err.Error())
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
	go func() {
		err = server.Serve(listen)
		if err != nil {
			panic("server启动失败")
		}
	}()
	// 程序优雅退出  对consul注销相关的service的实例
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	//timer := time.NewTimer(10 * time.Second)
	client.Agent().ServiceDeregister(id)
}
