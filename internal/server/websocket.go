package server

import (
        "encoding/json"
        
        //"github.com/fangelod/webrtc-test/internal/rtc"

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
        // Data RTCSessionDescription ?
        Action  string  `json:"action"`
        CallId  string  `json:"callId"`
        Name    string  `json:"name"`
        //Offer   webrtc.SessionDescription  `json:"offer"`

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

                        //case  "NEW_CALL" {
                        //        newCallAction(m)        
                        //}

                        //case "JOIN_CALL" {
                        //        joinCallAction(m)
                        //}

                        //case "LEAVE_CALL" {
                        //        leaveCallAction(m)
                        //}               
	        }
		conn.WriteMessage(t, msg)
        }
}

//func newCallAction(m message) {
//        answer, callId, err := rtc.NewCall(m.Offer, m.name, "")
//        if err != nil {
//                //messageBuilder with error
//        }
//
//        //messageBuilder()
//}

//func messageBuilder(){
//}
