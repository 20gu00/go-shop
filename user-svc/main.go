package main

import (
	"flag"
	"fmt"
	"go-shop/user-svc/model"
	"go-shop/user-svc/pb"
	"go-shop/user-svc/service"
	"google.golang.org/grpc"
	"net"
)

func main() {
	// 返回的是指针
	ip := flag.String("ip", "0.0.0.0", "ip")
	port := flag.String("port", "50051", "port")
	address := fmt.Sprintf("%s:%s", *ip, *port)

	flag.Parse()

	//userDao := model.NewUserDao()
	userDao := &model.User{}
	r, t, e := userDao.GetUserList()
	fmt.Println(r, t, e)
	// 创建server
	server := grpc.NewServer()

	// 注册user service
	pb.RegisterUserServer(server, &service.UserServer{})

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
