package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-web/common/initdo"
	"user-web/common/setUp/config"
	"user-web/common/tool"
	"user-web/router"
)

var c int = 1

func main() {
	var confFile string
	flag.StringVar(&confFile, "conf", "", "配置文件")
	flag.Parse()
	//读取配置文件,加载配置文件需要时间如果用goroutine方式去加载最好主goroutine阻塞一会,不然那拿到的配置值为空
	if err := config.ConfRead(confFile); err != nil {
		fmt.Printf("读取配置文件失败, err:%v\n", err)
		panic(err)
	}

	ch := make(chan int)
	// 不适用goroutine,主协程会阻塞,InitDo阻塞
	go func() {
		initdo.InitDO(ch)
	}()
	r := router.InitRouter()

	// 考虑一个问题,如果是本地开发环境端口固定也就是配置文件手动写定没问题,那如果服务很多,或者一台主机上很多服务,手动管理端口就不合适
	// 线上环境应该获取可用端口,至于服务暴露通过网关等,这些可以在负载均衡中做
	port, err := tool.GetFreePort()
	if err == nil {
		config.Conf.Port = port
	}

	server := http.Server{
		Addr:           fmt.Sprintf(":%d", config.Conf.Port),
		Handler:        r,
		ReadTimeout:    time.Duration(config.Conf.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.Conf.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << config.Conf.MaxHeader,
	}

	go func() {
		zap.L().Info("[Info]",
			zap.String("程序名称", viper.GetString("app_name")),
			zap.String("程序版本", viper.GetString("version")),
			zap.Int("server port", viper.GetInt("app_port")),
		)
		fmt.Println("[Info] server port:", viper.GetInt("app_port"))
		if err := server.ListenAndServe(); err != nil { //阻塞
			zap.L().Info("[Info] web server启动失败", zap.Error(err))
		}

	}()

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	ch <- c

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		zap.L().Fatal("server不正常退出,shutdown", zap.Error(err))
	}

	zap.L().Info("server退出了")

}
