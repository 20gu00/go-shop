package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"time"
	"user-web/common"
	user2 "user-web/model/user"
	"user-web/pb"
)

// 用户登录(api层做即可,调用依赖的rpc)
func UserPasswdLogin(ctx *gin.Context) {
	// 使用validate做参数规则教研,当然也可以自己写
	param := user2.UserPasswdLoginForm{}
	// json方式获取参数
	if err := ctx.ShouldBindJSON(&param); err != nil {
		// ctx 指针
		ValidateParam(ctx, err)
		return
	}
}

// 获取用户列表
func GetUserList(ctx *gin.Context) {
	// 方便开发,暂时直接写,viper已经做好了,也可以通过viper来做这些配置信息
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Error("[ GetUserList ]连接grpc server失败", zap.Error(err))
	}

	userClient := pb.NewUserClient(conn)

	// 获取参数 ShouldBindJSON  (json传参)
	// 设置query参数
	pNum, _ := strconv.Atoi(ctx.DefaultQuery("pnum", "0"))
	pSize, _ := strconv.Atoi(ctx.DefaultQuery("psize", "10"))

	res, err := userClient.GetUserList(context.Background(), &pb.PageInfo{
		PNum:  uint32(pNum),
		PSize: uint32(pSize),
	})
	if err != nil {
		// 也就是这个服务端的服务有错,客户端不应该异常的
		zap.L().Error("[ GetUserList ]查询用户失败")
		GrpcErrorToHttp(err, ctx)
	}

	// 时间戳转成time不用管时区
	//loc, _ := time.LoadLocation("Asia/Shanghai")

	resWeb := make([]interface{}, 0)
	for _, v := range res.UserListData {
		//data := make(map[string]interface{}, 0)
		user := user2.UserListRes{
			Id:       v.Id,
			NickName: v.Nickname,
			Gender:   v.Gender,
			Birthday: time.Unix(int64(v.Birthday), 0).Format("2006-01-02"),
			Mobile:   v.Mobile,
		}

		resWeb = append(resWeb, user)
		//data["id"] = v.Id
		//data["nickname"] = v.Nickname
		//data["gender"] = v.Gender
		//// 可以转换好再发送给前端或者前端在做转换
		//data["birthday"] = v.Birthday
		//data["mobile"] = v.Mobile
	}

	ctx.JSON(http.StatusOK, resWeb)
}

// grpc错误转换成http错误
func GrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		// 拿到grpc的错误
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "user grpc service不可用",
				})
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			// grpc的内部错误没必要过多返回
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": e.Code(),
					"msg":  "其他错误类型" + e.Message(),
				})
			}
			return
		}
	}
}

// 参数验证
func ValidateParam(ctx *gin.Context, err error) {
	err2, ok := err.(validator.ValidationErrors)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": err2.Error(),
		})
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"msg": err2.Translate(common.Trans),
	})
}
