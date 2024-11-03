package main

import (
	"fmt"
	"net"
	"time"
)

// checkPort 尝试连接指定主机和端口，返回端口是否开放
func checkPort(host string, port int, timeout time.Duration) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		fmt.Printf("端口 %d 关闭：%v\n", port, err)
		return false
	}
	conn.Close()
	fmt.Printf("端口 %d 开放\n", port)
	return true
}

func main() {
	host := "10.43.26.206"     // 要探测的主机
	port := 80                 // 要探测的端口
	timeout := 5 * time.Second // 超时时间设置为5秒

	// 尝试探测端口
	checkPort(host, port, timeout)
}
