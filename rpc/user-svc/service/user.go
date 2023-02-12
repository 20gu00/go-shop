package service

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
	"user-rpc/global/salt-passwd"
	"user-rpc/model"
	pb2 "user-rpc/pb"
)

// user grpc service
type UserServer struct {
	// 实现mustEmbedUnimplementedUserServer
	*pb2.UnimplementedUserServer
}

// 获取用户列表
func (u *UserServer) GetUserList(ctx context.Context, req *pb2.PageInfo) (rsp *pb2.UserListRes, err error) {
	//userDao := model.NewUserDao()
	userDao := model.NewUserDao()
	_, total, err := userDao.GetUserList()
	if err != nil {
		return nil, err
	}

	// 初始化
	res := &pb2.UserListRes{}
	res.Total = total

	// 拿到全部的user然后做分页,使用的gorm的分页
	userPaginate := userDao.Paginate(int(req.PNum), int(req.PSize))
	for _, user := range userPaginate {
		userInfoRes := Model2Res(user)
		// 使用切片要注意是否已经初始化也就是分配内存地址,不要使用rsp
		res.UserListData = append(res.UserListData, &userInfoRes)
	}
	return res, nil
}

func Model2Res(user model.User) pb2.UserInfo {
	userInfo := pb2.UserInfo{
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
func (u *UserServer) GetUserByMobile(ctx context.Context, req *pb2.Mobile) (rsp *pb2.UserInfo, err error) {
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
func (u *UserServer) GetUserById(ctx context.Context, req *pb2.Id) (res *pb2.UserInfo, err error) {
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

// 创建用户
func (u *UserServer) CreateUser(ctx context.Context, req *pb2.CreateUserInfo) (res *pb2.UserInfo, err error) {
	// 远程调用很多server的拦截器fmt.Println都不是server上输出了

	userDao := model.NewUserDao()
	user, result := userDao.GetUserByMobile(req.Mobile)
	if result.RowsAffected == 1 {
		// 服务端最好也打印或者日志一下,不然这里return,客户端那边报错退出,相关的错误信息较少
		fmt.Println("用户存在")
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}

	//fmt.Println("1", user.ID) // 1,0
	// 没有查询到,填充user
	user.NickName = req.Nickname
	user.Mobile = req.Mobile

	user.Password = salt_passwd.SaltPassword(req.Password)

	//fmt.Println("90", user.ID) //90 0
	// 创建用户
	// 注意create时,该结构体会拿到id,主建
	result = userDao.CreateUser(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	//fmt.Println("2", user.ID) //2,41  会拿到当前记录的id,即使创建的时候没有输入id,主建
	userInfo := Model2Res(user)
	//fmt.Println("3", userInfo.Id) //3,41
	return &userInfo, nil
}

// 用户个人中心更新用户
func (u *UserServer) UpdateUser(ctx context.Context, req *pb2.UpdateUserInfo) (*empty.Empty, error) {
	userDao := model.NewUserDao()
	user, result := userDao.GetUserById(req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	// 时间戳->time
	birthday := time.Unix(int64(req.Birthday), 0)
	user.Birthday = &birthday
	user.NickName = req.Nickname
	user.Gender = req.Gender

	result = userDao.UpdateUser(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	return &empty.Empty{}, nil
}

// 检查用户密码
func (u *UserServer) ValidatePassword(ctx context.Context, req *pb2.PasswordInfo) (*pb2.ValidateRes, error) {
	check := salt_passwd.ParsePassword(req.Password, req.EncryptedPassword)
	return &pb2.ValidateRes{
		Success: check,
	}, nil
}
