package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 根据需要配置允许的请求来源
	},
}

func GetWsMes(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to websocket", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// WebSocket 连接逻辑
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			return
		}
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			return
		}
	}
}
