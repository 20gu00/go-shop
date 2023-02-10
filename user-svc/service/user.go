package service

import (
	"context"
	"go-shop/user-svc/pb"
)

// user grpc service
type UserServer struct{}

// 获取用户列表
func (u *UserServer) GetUserList(context.Context, *pb.PageInfo) (*pb.UserListRes, error) {

}