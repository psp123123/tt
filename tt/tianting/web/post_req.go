package Web

import (
	"fmt"
	"tianting/comm"
	"tianting/conn"
	Conn "tianting/conn"
	_ "tianting/conn"

	"github.com/gin-gonic/gin"
)

var log comm.ConsoleLogger

type install_context struct {
	install_url     string
	install_name    string
	install_version string
}

// func InstallPost(c *gin.Context) {
// 	err := Conn.InitDB()
// 	if err != nil {
// 		log.Error("MySQL connected is failed!")
// 	}
// 	ret_ist := Conn.InsertHostData("1")
// 	log.Debug("Other Request'Tag is Ready! Reslut:%v\n", ret_ist)

// 	var req_context install_context
// 	req_context.install_url = c.PostForm("install_url")
// 	req_context.install_name = c.PostForm("install_name")
// 	req_context.install_version = c.PostForm("install_version")

// 	log.Debug("得到安装URL:%v,定义名称:%v,定义版本：%v\n", req_context.install_url, req_context.install_name, req_context.install_version)
// 	c.Request.Method = "GET"
// 	c.Request.URL.Path = "/index"
// }

func InsertHostPost(c *gin.Context) {
	H_host := c.PostForm("H_host")
	H_hostname := c.PostForm("H_hostname")
	H_core := c.PostForm("H_core")
	H_free := c.PostForm("H_free")
	H_disk := c.PostForm("H_disk")
	fmt.Printf("得到浏览器请求的IP信息：H_host=%vH_hostname=%vH_core=%vH_free=%vH_disk=%v", H_host, H_hostname, H_core, H_free, H_disk)
	ret_conn := Conn.InitDB()
	fmt.Printf("调用mysql信息：%v\n", ret_conn)
	ret := Conn.InsertHostData(H_host, H_hostname, H_core, H_free, H_disk)
	fmt.Println("mysql 数据库插入后的返回结果id", ret)
	c.Request.Method = "GET"
	c.Request.URL.Path = "/host_list_all"
}

func Submit(c *gin.Context) {
	type Configs struct {
		// 如果想要指定返回的json为其他别名，则可以使用`json:"username"`定义json格式的别名
		SoftwareName string `json:"softwareName" form:"softwareName"`
		Mode         string `json:"mode" form:"mode"`
		//Way          string `json:"way"`
		Version string `json:"version" form:"version"`
		Path    string `json:"path" form:"path" `
		//BuildArgs    string `json:"build_args"`
		Config string `json:"config" form:"config"`
	}
	var configBody Configs
	err := c.Bind(&configBody)
	// 判断json请求数据结构与定义的结构体有没有绑定成功
	log.Debug("------------------->得到的post表单%v", string(c.ContentType()))
	log.Debug("------------------->得到SoftName:%v", configBody.SoftwareName)
	log.Debug("------------------->得到SoftPath:%v", configBody.Path)
	c.Request.Method = "GET"
	c.Request.URL.Path = "/host_list_all"
	if err != nil {
		c.JSON(200, gin.H{
			"err_no":  400,
			"message": "Post Data Err",
		})
	} else {
		//install_res, rpm_file_path := fpm.Install(configBody.SoftwareName, configBody.Mode, configBody.Version, configBody.Path, configBody.Config)
		// c.JSON(http.StatusOK, gin.H{
		// 	"SoftwareName": configBody.SoftwareName,
		// 	"Mode":         configBody.Mode,
		// 	"源文件下载路径是":     install_res,
		// 	"制作完成的路径是":     rpm_file_path,
		// })
		c.Request.Method = "GET"
		c.Request.URL.Path = "192.168.56.11/#/page/table.html"
	}

}

func Heartbeat(c *gin.Context) {
	type DataStruct struct {
		Ip string `json:"ip" form:"ip"`
		Id string `json:"id" form:"id"`
	}
	var postData DataStruct
	err := c.Bind(&postData)
	if err != nil {
		c.JSON(200, gin.H{
			"err_no":  401,
			"message": "Client post error",
		})
	} else {
		result := conn.InitRedisConn(postData.Ip, "1")
		fmt.Printf("get redis func result is:%v\n", result)
		c.JSON(200, gin.H{
			"message": result,
		})
	}
	log.Debug("------------------->得到的post表单%v", string(c.ContentType()))
	log.Debug("------------------->得到IP:%T", postData.Ip)
	log.Debug("------------------->得到ID:%v", postData.Id)

}
