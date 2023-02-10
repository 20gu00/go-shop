package service

import (
	"context"
	"go-shop/user-svc/model"
	"go-shop/user-svc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// user grpc service
type UserServer struct{}

// 获取用户列表
func (u *UserServer) GetUserList(ctx context.Context, req *pb.PageInfo) (rsp *pb.UserListRes, err error) {
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
		userInfoRes := Model2Res(user)
		rsp.UserListData = append(rsp.UserListData, &userInfoRes)
	}
	return rsp, nil
}

func Model2Res(user model.User) pb.UserInfo {
	userInfo := pb.UserInfo{
		Id:       user.ID,
		Nickname: user.NickName,
		Password: user.Password,
		Gender:   user.Gender,
		Role:     int32(user.Role),
		// grpc的message的结构体的字段由默认值,不能随意将nil赋值进去,序列化的时候会异常
		// 搞清楚哪些字段是由默认值,birthday一开始可能为空
		// Birthday: user.Birthday
	}
	if user.Birthday != nil {
		// 时间戳
		userInfo.Birthday = uint32(user.Birthday.Unix())
	}
	return userInfo
}

// 通过mobile查询用户
func GetUserByMobile(ctx context.Context, req *pb.Mobile) (rsp *pb.UserInfo, err error) {
	userDao := model.NewUserDao()
	user, result := userDao.GetUserByMobile(req.Mobile)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户查找不到")
	}

	// 其他错误,比如数据库连接不上等
	if result.Error != nil {
		return nil, result.Error
	}

	userInfo := Model2Res(user)
	return &userInfo, nil
}

// 通过id查找用户
func GetUserById(ctx context.Context, req *pb.Id) (res *pb.UserInfo, error) {
	userDao := model.NewUserDao()
	user, result := userDao.GetUserById(req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户查找不到")
	}

	// 其他错误,比如数据库连接不上等
	if result.Error != nil {
		return nil, result.Error
	}

	userInfo := Model2Res(user)
	return &userInfo, nil
}
