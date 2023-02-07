package Web

import (
	"net/http"
	"strconv"
	"tianting/fpm"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// var log comm.ConsoleLogger
var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func GetWsMes(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Debug("upgrade:", err)
		return
	}
	defer ws.Close()

	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Debug("read:", err)
			break
		}
		log.Debug("recv: %s", message)

		Num := 0
		for {
			Num = Num + 1
			EndStr := "start"
			if Num > 10 {
				EndStr = "End"
			}
			//需要返回给浏览器的数据ResPonseMes,后续数据类型为[]byte(string(message) + " " + strconv.Itoa(count))
			//1. 下载判断值并返回给浏览器
			fpm.Install
			//如果值为End，则数据发送完毕,调用函数返回需要对end字符做处理
			ResPonseMes := []byte(strconv.Itoa(Num))
			log.Debug("数据更新为%v", Num)
			err = ws.WriteMessage(mt, ResPonseMes)
			if err != nil {
				log.Debug("write:", err)
				ws.Close()
			}
			//log.Debug("reponse data is %s----%v", ResPonseMes, mt)
			time.Sleep(time.Second)
			if string(EndStr) == "End" {
				log.Debug("end string is %v", EndStr)

				ws.Close()
				break

			}
		}

	}
}
