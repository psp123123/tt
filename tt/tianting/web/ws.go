package Web

import (
	"encoding/json"
	"net/http"
	Conn "tianting/conn"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Host represents the structure of the host data
type Host struct {
	ID        int    `gorm:"column:id"`
	HHost     string `gorm:"column:H_host"`
	HHostname string `gorm:"column:H_hostname"`
	HCore     int    `gorm:"column:H_core"`
	HFree     int    `gorm:"column:H_free"`
	HDisk     int    `gorm:"column:H_disk"`
	HStatus   int    `gorm:"column:H_status"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (adjust this for production)
	},
}

// ClientRequest represents the structure of the message received from the client
type ClientRequest struct {
	IDs []string `json:"ids"` // Array of IDs as expected from the client message
}

// GetWsMes handles WebSocket connections, receives queries, and continuously sends host data
func GetWsMes(c *gin.Context) {
	log.Debug("WebSocket 请求接收中...")

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Error("连接升级失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "连接升级失败"})
		return
	}
	defer func() {
		log.Debug("关闭 WebSocket 连接")
		conn.Close()
	}()

	log.Debug("WebSocket 连接已建立")

	// Connect to the database
	db, err := gorm.Open(mysql.Open("root:123456@tcp(10.43.26.206:3306)/auto_deploy"), &gorm.Config{})
	if err != nil {
		log.Error("数据库连接失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "连接数据库失败"})
		return
	}
	log.Debug("数据库连接成功")

	var clientReq ClientRequest
	// 读取客户端的初始消息
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Error("读取客户端初始消息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取客户端消息失败"})
		return
	}
	if err := json.Unmarshal(message, &clientReq); err != nil {
		log.Error("解析客户端初始消息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析客户端消息失败"})
		return
	}
	//log.Debug("客户端请求查询主机 ID 列表: %v", clientReq.IDs)

	// 定期查询数据库并发送数据
	queryInterval := time.NewTicker(1 * time.Second)
	defer queryInterval.Stop()

	for {
		select {
		case <-queryInterval.C:
			if err := queryDatabaseAndSendHostData(conn, db, c, clientReq.IDs); err != nil {
				log.Error("数据库查询或数据发送失败: %v", err)
				return
			}
		}
	}
}

func queryDatabaseAndSendHostData(conn *websocket.Conn, db *gorm.DB, c *gin.Context, ids []string) error {
	var hosts []Host
	query := db

	// 清空之前的状态集合
	StatusSetToCli := []map[string]string{}

	// If the client provided a list of IDs, filter by those IDs
	if len(ids) > 0 {
		query = query.Where("id IN ?", ids)
	}

	// Execute the query and check for errors
	//log.Debug("查询数据库中主机数据, 主机 ID 列表: %v - %T", ids, ids)
	if err := query.Find(&hosts).Error; err != nil {
		log.Error("查询数据库错误: %v", err)
		return err
	}

	//log.Debug("查询到 %d 个主机记录", len(hosts))

	for _, HID := range ids {
		//log.Debug("---Get id is: %v", HID)
		resStatus := Conn.GetHostContext(HID)
		StatusSetToCli = append(StatusSetToCli, map[string]string{"id": HID, "host": resStatus["H_host"], "H_status": resStatus["H_status"]})
		//log.Debug("查询数据主机状态结果：%v", resStatus)
		//log.Debug("组装好的数据即将发给客户端:%v", StatusSetToCli)
	}

	// Prepare and send the host data to the client
	msg, err := json.Marshal(StatusSetToCli)
	if err != nil {
		log.Error("序列化主机数据失败: %v", err)
		c.Error(err)
		return err
	}

	// Send the message to the WebSocket connection
	//log.Debug("发送主机数据到客户端: %v", string(msg))
	if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Error("发送消息到客户端失败: %v", err)
		c.Error(err)
		return err
	}

	return nil
}
