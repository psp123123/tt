package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return "ERR"
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println(ipnet.IP.String())
				return ipnet.IP.String()
			}
		}
	}
	return "ERR"
}

const Url = "http://47.93.37.151/heartbeat"

func PostInformation() {
	LocalIp := GetLocalIP()

	resp, err := http.PostForm(Url,
		url.Values{
			"ip": {LocalIp},
			"id": {"10"},
		},
	)
	if err != nil {
		fmt.Printf("post error :%v\n", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("get resp failed, err:%v\n", err)
		return
	}
	fmt.Printf("get message:%v\n", string(b))
}

func main() {
	for {
		PostInformation()
		time.Sleep(time.Second)
	}
}
