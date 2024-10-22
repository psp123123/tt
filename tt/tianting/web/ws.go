package Web

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Host represents the structure of the host data
type Host struct {
	H_id       int    `gorm:"column:H_id"`
	H_host     string `gorm:"column:H_host"`
	H_hostname string `gorm:"column:H_hostname"`
	H_core     int    `gorm:"column:H_core"`
	H_free     int    `gorm:"column:H_free"`
	H_disk     int    `gorm:"column:H_disk"`
	H_status   int    `gorm:"column:H_status"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (adjust this for production)
	},
}

// GetWsMes handles WebSocket connections and continuously sends host data
func GetWsMes(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// Log the error and return a response
		c.JSON(http.StatusInternalServerError, gin.H{"error": "连接升级失败"})
		return
	}
	defer conn.Close() // Ensure the connection is closed when the function exits

	// Connect to the database
	db, err := gorm.Open(mysql.Open("root:123456@tcp(10.43.26.206:3306)/auto_deploy"), &gorm.Config{})
	if err != nil {
		// Log the error and return a response
		c.JSON(http.StatusInternalServerError, gin.H{"error": "连接数据库失败"})
		return
	}

	// Continuously send host data
	go func() {
		for {
			var hosts []Host
			if err := db.Find(&hosts).Error; err != nil {
				// If an error occurs while querying the database, log it and return
				c.Error(err) // Log the error for debugging
				return
			}

			// Prepare and send the host data
			for _, host := range hosts {
				msg, err := json.Marshal(host)
				if err != nil {
					c.Error(err) // Log JSON serialization errors
					return
				}

				// Send the message to the WebSocket connection
				if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					c.Error(err) // Log send message errors
					return
				}
			}

			time.Sleep(1 * time.Second) // Send data every second
		}
	}()
}
