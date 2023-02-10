package main

import (
	"context"
	"fmt"
	"go-shop/user-svc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/*
	测试user grpc service
	测试userList和validatePasswd接口
*/

var (
	userClient pb.UserClient
	conn       *grpc.ClientConn
)

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	//defer conn.Close()
	// 调用的函数中关闭

	userClient = pb.NewUserClient(conn)
}

func TestGetUserList() {
	res, err := userClient.GetUserList(context.Background(), &pb.PageInfo{
		PNum:  1,
		PSize: 10, // 获取10条数据
	})
	if err != nil {
		panic(err)
	}

	for _, user := range res.UserListData {
		fmt.Println(user.Nickname, user.Password, user.Mobile)
		checkRes, err := userClient.ValidatePassword(context.Background(), &pb.PasswordInfo{
			// 原始密码,也就是用户输入的密码
			Password: "admin12345",
			// 加密后的密码
			EncryptedPassword: user.Password,
		})
		if err != nil {
			panic(err)
		}

		fmt.Println(checkRes)
	}
}

func main() {
	// 前提是grpc server启动
	Init()
	TestGetUserList()
	defer conn.Close()
}
