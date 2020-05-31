package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader websocket.Upgrader

func websocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		conn.WriteMessage(t, msg)
	}
}
