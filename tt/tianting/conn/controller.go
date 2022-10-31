package conn

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)



func Controller(c *gin.Context)  {
		host_context := c.PostForm("host_client_context")
		host_status := c.PostForm("host_client_status")
		if string(host_status) == "server_on" {
			fmt.Println("client is connected")
		}
		fmt.Println(host_context,host_status)
		c.String(http.StatusOK,"Server_Response")
	}

func ControllerPlus(c *gin.Context) {
	//c.String(http.StatusOK,"you will be installed sth!!!!!!!!!!")
	c.JSON(http.StatusOK,gin.H{
		"get_url": "https://mirrors.tuna.tsinghua.edu.cn/centos/timestamp.txt",
		"service": "redis",
		"version": "v3.8",
	})
}

func JudgePing(c *gin.Context) {
	InitDB()
	ret_one := GetTagOne(1)
	fmt.Printf("获得数据库是否安装服务接口： %v （0：不安装，1：安装）\n",ret_one)
	ret_Tag,err := strconv.Atoi(ret_one)
	if err != nil {
		fmt.Printf("get tag is error:%v\n",err)
	}
	if ret_Tag == 0 {
		fmt.Printf("收到客户端心跳信息： %v\n",ret_Tag)
		Controller(c)
	}else {
		fmt.Printf("get ret_Tag is %v\n",ret_Tag)
		ControllerPlus(c)
	}
}

