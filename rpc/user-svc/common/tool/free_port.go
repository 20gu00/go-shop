package tool

import "net"

// 获取可用的端口,尤其是多个服务在同一台主机上

func GetFreePort() (int, error) {
	address, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, nil
	}
	l, err := net.ListenTCP("tcp", address)
	if err != nil {
		return 0, nil
	}
	defer l.Close()

	// 返回一个可用的端口号
	return l.Addr().(*net.TCPAddr).Port, nil
}
