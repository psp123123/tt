package Web

import (
	"fmt"
	"io/ioutil"
	"net/http"
	Conn "tianting/conn"
	"tianting/tools"

	"github.com/gin-gonic/gin"
)

func IndexGet(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"title": "我是测试", "ce": "123456"})
}

func EditGet(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"edit.html",
		gin.H{"status": "ok"},
	)
}
func InstallGet(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"install.html",
		gin.H{"status": "ok"},
	)
}

// 取svc_name表内容
func Snames(c *gin.Context) {
	Conn.InitDB()
	Get_ret, ter := Conn.GetDataAll_svc_name()
	println(Get_ret)
	c.JSON(http.StatusOK, gin.H{"allData": Get_ret, "category": ter})

}

type ReqHostList struct {
	// Host   string `binding:"required"`
	// Status string `binding:"required`
	Hid int `binding:"required"`
}

// 获取已存在的机器列表信息--单行
func Host_list(c *gin.Context) {
	Conn.InitDB()
	var reqhostlist ReqHostList
	if err := c.ShouldBind(&reqhostlist); err != nil {
		c.String(500, fmt.Sprint(err))
		return
	}
	Get_ret := Conn.GetHostContext(c.Query("Hid"))
	fmt.Printf("get ret is -------%v", Get_ret)

	c.String(200, fmt.Sprintf("参数绑定后:----%v--获取的结果:主机:%v--是否在线[%v]\n", reqhostlist, Get_ret["H_host"], Get_ret["H_status"]))
	c.JSON(http.StatusOK, gin.H{

		"host":   Get_ret["H_host"],
		"status": Get_ret["H_status"],
	})

}

// 获取已存在的机器列表信息--多行
func Host_list_all(c *gin.Context) {
	Conn.InitDB()

	Get_ret := Conn.GetHostContextAll()

	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"msg":   "",
		"count": 1000,
		// "H_host":   "192.168.100.11",
		// "H_status": "1",
		"data": Get_ret,
	})

}

// 获取配置文件，得到数据：{"server_name":"redis",“server_configs”:"各个配置项"}
func Configs(c *gin.Context) {
	Conn.InitDB()
	server_name := c.DefaultQuery("server_name", "使用接口方式是，/server_configs?server_name=redis")
	fmt.Printf("从请求URI中得到的服务名参数是：%v\n", server_name)
	Getconfg_name := Conn.GetServerConfigName(server_name)
	//判断uri传入的参数是否可在数据库中查到，后续结果都基于这个判断
	if Getconfg_name == "MysqlErr" {
		fmt.Printf("因为数据库返回错误，将返回0值结果：%v\n", Getconfg_name)
		c.JSON(http.StatusNotFound, gin.H{server_name: Getconfg_name, "server_content": "传入参数为空，请检查"})
	} else {
		fmt.Printf("从Mysql数据库中configs函数得到的服务配置文件的名称：%v\n", Getconfg_name)
		//var base_dir string
		//sysType:=runtime.GOOS
		//if sysType == "linux" {
		base_dir := tools.JudgeOS()
		file_path := base_dir + "src/tianting/instance/" + Getconfg_name
		fmt.Printf("=====================>>>获取文件全路径：%v\n", file_path)
		content, err := ioutil.ReadFile(file_path)
		if err != nil {
			fmt.Println("read file failed, err:", err)
			c.JSON(http.StatusForbidden, gin.H{server_name: Getconfg_name, "报错原因：": server_name + "所需要的配置文件Config file Not found!"})
		} else {
			c.JSON(http.StatusOK, gin.H{server_name: Getconfg_name, "server_content": string(content)})
		}

	}

}

// 删除host单行数据
func Del_host_one(c *gin.Context) {
	con := c.Query("H_id")
	n := Conn.DelHostData(con)
	fmt.Printf("get client id %v\n", con)
	fmt.Printf("删除后的结果是%v\n", n)
	if n > 0 {
		c.String(http.StatusOK, "200")
	} else {
		c.String(http.StatusOK, "500")
	}
}

func Test(c *gin.Context) {

	config_path := "./src/tianting/instance/"
	fmt.Printf("=======================>>>获取文件的父目录：%v\n", config_path)

	ret := Conn.GetServerConfigName("redis")
	fmt.Printf("Test函数得到的值是%v\n", ret)
	//指定默认值

	//http://localhost:8080/user 才会打印出来默认的值
	name := c.DefaultQuery("name", "枯藤")
	file_all_path := config_path + name + ".conf"
	fmt.Printf("=====================================>>>得到文件全路径：%v\n", file_all_path)
	c.String(http.StatusOK, fmt.Sprintf("hello %s", name))

}
