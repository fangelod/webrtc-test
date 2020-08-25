package server

import (
	"encoding/json"
	"net/http"

	"github.com/fangelod/webrtc-test/internal/rtc"
	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

func newCallHandler(c *gin.Context) {
	sdp := rtc.RTCSessionDescription{}
	if err := json.NewDecoder(c.Request.Body).Decode(&sdp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer, callId, err := rtc.NewCall(sdp.Description, sdp.Name, sdp.User)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"callId": callId, "sdp": answer})
}

func getCallsHandler(c *gin.Context) {
	calls, err := rtc.GetCalls()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, calls)
}

func joinCallHandler(c *gin.Context) {
	idString := c.Param("id")
	id, err := uuid.FromString(idString)
	if err != nil {
	        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

        sdp := rtc.RTCSessionDescription{}
        if err := json.NewDecoder(c.Request.Body).Decode(&sdp); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        answer, err := rtc.JoinCall(sdp.User, sdp.Description, id)
        if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

        c.JSON(http.StatusOK, gin.H{"sdp": answer})
}

func leaveCallHandler(c *gin.Context) {
	idString := c.Param("id")
	id, err := uuid.FromString(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := rtc.LeaveCall("", id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, "")
}

func renegotiateCallHandler(c *gin.Context) {
	idString := c.Param("id")
	id, err := uuid.FromString(idString)
	if err != nil {
	        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

        sdp := rtc.RTCSessionDescription{}
        if err := json.NewDecoder(c.Request.Body).Decode(&sdp); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        answer, err := rtc.RenegotiateCall(sdp.User, id)
        if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

        c.JSON(http.StatusOK, gin.H{"sdp": answer})
}

func iceCandidateHandler(c *gin.Context) {
	idString := c.Param("id")
	id, err := uuid.FromString(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	can := rtc.CandidateFromClient{}
	if err := json.NewDecoder(c.Request.Body).Decode(&can); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := rtc.AddICECandidateForUser(id, can.User, can.Candidate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, "")
}

func iceServersHandler(c *gin.Context) {
	servers := rtc.GetIceServers()

	c.JSON(http.StatusOK, servers)
}
