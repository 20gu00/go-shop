package global

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"user-web/common/setUp/config"
	"user-web/pb"
)

// 全局的user rpc服务的client

var (
	GlobalUserClient     pb.UserClient
	UserClientFromConsul pb.UserClient
)

func CreateUserClientFromConsul() {
	c, err := grpc.Dial(
		// consul中的service_name
		// wait等待解析多长时间
		// limit有多个服务只要多少个
		// tag=manual过滤作用
		// https://github.com/mbobakov/grpc-consul-resolver

		fmt.Sprintf("consul://%s:%d/%s?wait=15s", config.Conf.ConsulConfig.Host, config.Conf.ConsulConfig.Port, config.Conf.UserRpcConfig.Name),
		// grpc.WithInsecure(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),

		// 负载均衡算法
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// 创建一个grpc client,操作grpc server,这事是直接操作consul返回的信息
	UserClientFromConsul = pb.NewUserClient(c)
	//可以使用这个client调用方法了
}
