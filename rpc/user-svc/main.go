package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"

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
	ip := flag.String("ip", "0.0.0.0", "ip")
	port := flag.String("port", "50051", "port")
	address := fmt.Sprintf("%s:%s", *ip, *port)

	flag.Parse()

	//读取配置文件,加载配置文件需要时间如果用goroutine方式去加载最好主goroutine阻塞一会,不然那拿到的配置值为空
	if err := config.ConfRead(confFile); err != nil {
		fmt.Printf("读取配置文件失败, err:%v\n", err)
		panic(err)
	}

	ch := make(chan int)
	go func() {
		global.Init(ch)
	}()

	// 创建server
	server := grpc.NewServer()

	// 注册user service
	pb.RegisterUserServer(server, &service.UserServer{})
	// 注册健康检查
	srv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, srv)

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
