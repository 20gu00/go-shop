package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"time"
	"user-web/common"
	"user-web/common/jwt"
	"user-web/common/setUp/config"
	"user-web/dao/redis"
	user2 "user-web/model/user"
	"user-web/pb"
)

var userClient pb.UserClient

func init() {
	// 方便开发,暂时直接写,viper已经做好了,也可以通过viper来做这些配置信息
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Error("连接grpc server失败", zap.Error(err))
	}

	userClient = pb.NewUserClient(conn)
}

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

	// 先进行验证码验证逻辑(机器人身份等)
	// 验证码id 验证码 是否清理掉(就是填写了验证码进行了一次访问,就不能用同一个验证码在次访问,要刷新)
	// postman访问时就是带上captcha接口的验证码id和验证码(前端根据captcha接口得到验证码id和验证码,然后提供验证码的空格让用户输入)
	if !captchaStore.Verify(param.CaptchaId, param.Captcha, true) {
		// 验证失败
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码错误",
		})
		return
	}
	// logic
	res, err := userClient.GetUserByMobile(context.Background(), &pb.Mobile{
		Mobile: param.Mobile,
	})
	if err != nil {
		err2, ok := status.FromError(err)
		if ok {
			switch err2.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "用户不存在",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "登录失败",
				})
			}
			return
		} else {
			// 用户存在,检查密码是否正确的
			if res2, err2 := userClient.ValidatePassword(context.Background(), &pb.PasswordInfo{
				// 原始密码
				Password: param.Password,
				// 加密密码
				EncryptedPassword: res.Password,
			}); err2 != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "登录失败",
				})
			} else {
				if res2.Success {
					/* cookie和session机制实现登录状态
						   流程:
						   		1.浏览器访问请求服务端,登录
								2.服务端查询数据库中的用户
								3.数据库返回用户信息给服务端
								4.服务端针对这个用户创建一个session和sessionid会保存到数据库中
								5.将sessionid返回给浏览器,设置到cookie中
								6.后续浏览器请求都会带上这个sessionid
								7.服务端通过sessionid确定用户状态和从数据库中获取用户相关信息
						   在微服务中的问题:
						其实sessionid其实就是用来让服务端获取出相应用户的信息
						比如浏览器从用户微服务中获取sessionid,然后它带着这个sessionid去访问商品微服务,而商品微服务的数据哭中没用用户信息,那么它应该去用户微服务的数据库中获取,单微服务应该是独立的,或者你可以用一个公用的数据库
						比如用redis而不是mysql来存放session信息,用作公用数据库,这就是分布式session,那么这个redis就要去扛住高并发

					json web token在微服务中更好用
						两点功能:身份验证和信息交换
						加密的jwt字符串
						各个微服务服务端加密 解密,浏览器不能解密
						使用同样的key加密解密即可
						token加解密来判断信息即可,不用像session一样在本地存储
						token一般放在header中

					*/

					///生成token
					token, err := jwt.GenToken(uint(res.Id), uint(res.Role), res.Nickname)
					if err != nil {
						ctx.JSON(http.StatusInternalServerError, gin.H{
							"msg": "token生成失败",
						})
						return
					}
					ctx.JSON(http.StatusOK, gin.H{
						"msg":      "登陆成功",
						"id":       res.Id,
						"token":    token, // 主要还是通过解析token来获取信息
						"Nickname": res.Nickname,
					})
				} else {
					ctx.JSON(http.StatusOK, gin.H{
						"msg": "登录失败",
					})
				}
			}
		}
	}
}

// 获取用户列表
func GetUserList(ctx *gin.Context) {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", config.Conf.ConsulConfig.Host, config.Conf.ConsulConfig.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		zap.L().Error("[ GetUserList ]创建consul的client失败")
		return
	}

	userRpcHost := ""
	userRpcPort := 0
	//data,err:=client.Agent().ServicesWithFilter(`Service == "user-rpc"`)
	// 这个格式很重要  或者转义比如\"
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service =="%s"`, config.Conf.UserRpcConfig.Name))
	if err != nil {
		zap.L().Error("[ GetUserList ]从consul总过滤服务失败")
		return
	}
	for _, v := range data {
		userRpcHost = v.Address
		userRpcPort = v.Port
		// 获取这个service任意一个负载即可
		break
	}

	if userRpcHost == "" {
		zap.L().Error("[ GetUserList ]获取rpc服务负载实例失败")
		return
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", userRpcHost, userRpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Error("连接grpc server失败", zap.Error(err))
	}

	userClient2 := pb.NewUserClient(conn)

	// 解析token获取相应的用户信息
	claim, _ := ctx.Get("claim")
	user, ok := claim.(jwt.MyClaims)
	if !ok {
		zap.L().Error("context的claim断言失败")
	}
	zap.L().Info("访问的用户: ", zap.Int("userId", int(user.ID)))

	// 获取参数 ShouldBindJSON  (json传参)
	// 设置query参数
	pNum, _ := strconv.Atoi(ctx.DefaultQuery("pnum", "0"))
	pSize, _ := strconv.Atoi(ctx.DefaultQuery("psize", "10"))

	res, err := userClient2.GetUserList(context.Background(), &pb.PageInfo{
		PNum:  uint32(pNum),
		PSize: uint32(pSize),
	})
	if err != nil {
		// 也就是这个服务端的服务有错,客户端不应该异常的
		zap.L().Error("[ GetUserList ]查询用户失败")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "查询用户失败",
		})
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

// 用户注册逻辑
func RegisterUser(ctx *gin.Context) {
	param := user2.UserRegisterInput{}
	// json方式获取参数
	if err := ctx.ShouldBindJSON(&param); err != nil {
		// ctx 指针
		ValidateParam(ctx, err)
		return
	}

	// 验证码校验
	// 验证码生成接口已经将验证码保存进redis
	v, err := redis.GetSmsCode(param.Mobile)
	if err != nil {
		zap.L().Error("通过手机号码从redis获取验证码失败")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err,
		})
	} else {
		if v != param.SmsCode {
			zap.L().Error("验证码校验错误")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error() + "验证码校验错误",
			})
		}
		return
	}

	// 注册用户
	user, err := userClient.CreateUser(context.Background(), &pb.CreateUserInfo{
		Nickname: param.Mobile,
		Mobile:   param.Mobile,
		Password: param.Password,
	})
	if err != nil {
		zap.S().Errorf("[ RegisterUser ]新建用户失败: %s", ctx)
		GrpcErrorToHttp(err, ctx)
		return
	}

	// 如果是注册即可登录,设置token
	token, err := jwt.GenToken(uint(user.Id), uint(user.Role), user.Nickname)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "token生成失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":      "登陆成功",
		"id":       user.Id,
		"token":    token, // 主要还是通过解析token来获取信息
		"Nickname": user.Nickname,
	})

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
