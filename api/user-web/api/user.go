package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"net/http"
	"user-web/pb"
)

func GetUserList(ctx *gin.Context) {
	// 方便开发,暂时直接写,viper已经做好了,也可以通过viper来做这些配置信息
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Error("[ GetUserList ]连接grpc server失败", zap.Error(err))
	}

	userClient := pb.NewUserClient(conn)

	res, err := userClient.GetUserList(context.Background(), &pb.PageInfo{
		PNum:  1,
		PSize: 10,
	})
	if err != nil {
		// 也就是这个服务端的服务有错,客户端不应该异常的
		zap.L().Error("[ GetUserList ]查询用户失败")
		GrpcErrorToHttp(err, ctx)
	}

	resWeb := make([]interface{}, 0)
	for _, v := range res.UserListData {
		data := make(map[string]interface{}, 0)
		data["id"] = v.Id
		data["nickname"] = v.Nickname
		data["gender"] = v.Gender
		// 可以转换好再发送给前端或者前端在做转换
		data["birthday"] = v.Birthday
		data["mobile"] = v.Mobile
		resWeb = append(resWeb, data)
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
