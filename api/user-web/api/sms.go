package api

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"user-web/dao/redis"
	user2 "user-web/model/user"
)

/*
	阿里云官网->短信服务->短信控制台->开通服务->国内短信->添加签名->验证码
	->返回列表->模板管理->添加模板->验证码->您正在进行手机注册，验证码为：${code}，5分钟内有效
	->返回列表->获取模板code

*/

// 使用第三方阿里云来发送短信
func Sms(ctx *gin.Context) {
	param := user2.SmsInput{}
	// json方式获取参数
	if err := ctx.ShouldBindJSON(&param); err != nil {
		// ctx 指针
		// ctx.Json()
		ValidateParam(ctx, err)
		return
	}
	//工作台->账号头像->AccessKey管理->AccessKeyId和AccessKeySecret
	client, err := dysmsapi.NewClientWithAccessKey("cn-beijing", "AccessKeyId", "AccessKeySecret")
	mobile := param.Mobile
	smsCode := GenSmsCode(5)
	if err != nil {
		//panic(err)
		return
	}
	r := requests.NewCommonRequest()
	r.Method = "POST"
	r.Scheme = "https" //http
	r.Domain = "dysmsapi.aliyuncs.com"
	r.Version = "2017-05-25"
	r.ApiName = "SendSms"
	r.QueryParams["RegionId"] = "cn-beijing"
	r.QueryParams["PhoneNumbers"] = mobile // 手机号码
	r.QueryParams["SignName"] = "xxx"      // 短信服务的签名,国内短信列表中的签名名称
	r.QueryParams["TemplateCode"] = "xxx"  // 模板的code
	r.QueryParams["TemplateParam"] = "{\"code\":" + smsCode + "}"
	response, err := client.ProcessCommonRequest(r)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("response is %#v\n", response)

	// 根据短信验证码注册

	// 保存验证码 手机号码key 验证码value 保存到redis中
	redis.InSmsCode(mobile, smsCode)
	// 返回给前端
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "短信验证码发送成功",
	})
}

// 生成验证码
func GenSmsCode(width int) string {
	num := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(num)
	rand.Seed(time.Now().Unix())
	var s strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&s, "%d", num[rand.Intn(r)])
	}
	return s.String()
}
