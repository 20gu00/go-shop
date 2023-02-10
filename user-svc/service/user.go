package service

import (
	"context"
	"go-shop/user-svc/model"
	"go-shop/user-svc/pb"
)

// user grpc service
type UserServer struct{}

// 获取用户列表
func (u *UserServer) GetUserList(ctx context.Context, req *pb.PageInfo) (*pb.UserListRes, error) {
	userDao := model.NewUserDao()
	userList, total, err := userDao.GetUserList()
	if err != nil {
		return nil, err
	}

	res := &pb.UserListRes{}
	res.Total = total

	// 拿到全部的user然后做分页,使用的gorm的分页
	userPaginate := userDao.Paginate(int(req.PNum), int(req.PSize))
	for _, user := range userPaginate {

	}
}

func Model2Res(user model.User) pb.UserInfo {
	userInfo := pb.UserInfo{
		Id:       user.ID,
		Nickname: user.NickName,
		Password: user.Password,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}

}