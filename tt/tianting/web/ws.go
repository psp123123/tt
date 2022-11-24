package Web

import (
	"net/http"
	"strconv"
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
			message := []byte(string(message) + " " + strconv.Itoa(Num))
			log.Debug("数据更新为%v", Num)
			err = ws.WriteMessage(mt, message)
			if err != nil {
				log.Debug("write:", err)
				break
			}
			log.Debug("reponse data is %s----%v", message, mt)
			time.Sleep(time.Second)
		}

	}
}
