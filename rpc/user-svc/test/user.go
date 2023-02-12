package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb2 "user-rpc/pb"
)

/*
	测试user grpc service
	测试userList和validatePasswd接口
*/

var (
	userClient pb2.UserClient
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

	userClient = pb2.NewUserClient(conn)
}

func TestGetUserList() {
	res, err := userClient.GetUserList(context.Background(), &pb2.PageInfo{
		PNum:  1,
		PSize: 2, // 获取10条数据
	})
	if err != nil {
		panic(err)
	}

	for _, user := range res.UserListData {
		fmt.Println(user.Nickname, user.Password, user.Mobile)
		checkRes, err := userClient.ValidatePassword(context.Background(), &pb2.PasswordInfo{
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

func TestCreateUser() {
	//options := &password.Options{16, 100, 30, sha512.New}
	//salt, encodedPwd := password.Encode("admin12345", options)
	//newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	for i := 200; i < 205; i++ {
		rsp, err := userClient.CreateUser(context.Background(), &pb2.CreateUserInfo{
			Nickname: fmt.Sprintf("用户%d", i),
			Mobile:   fmt.Sprintf("12345678%d", i),
			Password: "admin12345", // 传递明文密码,由服务端加密设置了 newPassword,
			Gender:   "Male",
		})
		// 要处理server返回的错误
		if err != nil {
			panic(err)
		}
		fmt.Println(rsp.Id)
	}

}

func main() {
	// 前提是grpc server启动
	Init()
	//TestGetUserList()
	TestCreateUser()
	//defer conn.Close()
	conn.Close()
}
