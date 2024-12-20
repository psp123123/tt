package main

import (
	Conn "tianting/conn"
	Web "tianting/web"

	"github.com/gin-gonic/gin"
)

func main() {
	var r = gin.Default()
	r.Static("/static", "./static/")
	//r.LoadHTMLGlob("tainting/tem/*")

	r.GET("/index", Web.IndexGet)
	//r.GET("/status", Web.StatusGet)
	r.GET("/edit", Web.EditGet)
	r.GET("/install", Web.InstallGet)
	r.GET("/server_list", Web.Snames)
	r.GET("/server_configs", Web.Configs)
	r.GET("/test", Web.Test)
	r.GET("/host_list", Web.Host_list)
	r.GET("/host_list_all", Web.Host_list_all)
	r.GET("/del_host_one", Web.Del_host_one)
	//ws协议
	r.GET("/ws", Web.GetWsMes)
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "page/wspush.html", nil)
	})

	//r.POST("/install", Web.InstallPost)
	r.POST("/post_context", Conn.JudgePing)
	r.POST("/insert_one", Web.InsertHostPost)
	r.POST("/del_host_one", Web.Del_host_one)
	r.POST("/heartbeat", Web.Heartbeat)

	r.POST("/submit", Web.Submit)
	//r.GET("/ws",Web.WebSocket)

	r.Run("0.0.0.0:3000")
}
