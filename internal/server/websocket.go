package server

import (
        "encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader websocket.Upgrader

	conns map[string]*websocket.Conn
)

func init() {
	conns = make(map[string]*websocket.Conn)
}


type message struct {
        //Need action for switch case
        //Need callId, username, sdp offer


}

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

	        m := message{}
	        if err := json.Unmarshal(msg, &m); err == nil {
                        //Switch case of actions
		        //if m.Action == "SAVE_NAME" {
		        //	conns[m.Data] = conn
		        //}
	        }

		conn.WriteMessage(t, msg)
	}
}
