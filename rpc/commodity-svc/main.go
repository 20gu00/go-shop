package main

import (
	"commodity-rpc/common/initDo"
	"commodity-rpc/common/setUp/config"
	"flag"
	"fmt"
)

func main() {
	var confFile string
	flag.StringVar(&confFile, "conf", "", "配置文件")

	if err := config.ConfRead(confFile); err != nil {
		fmt.Printf("读取配置文件失败, err:%v\n", err)
		panic(err)
	}
	initDo.Init()
}
