package server

import (
	"os"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

const (
	portEnvString    = "WEBRTC_TEST_PORT"
	webRootEnvString = "WEBRTC_TEST_WEBROOT"

	defaultPort    = ":8080"
	defaultWebRoot = "/opt/webrtc-test/public"
)

var (
	port    string
	webRoot string

	server *gin.Engine
)

func init() {
	envPort := os.Getenv(portEnvString)
	if envPort == "" {
		if err := os.Setenv(portEnvString, defaultPort); err != nil {
			panic(err)
		}
		port = defaultPort
	} else {
		port = envPort
	}

	envWebRoot := os.Getenv(webRootEnvString)
	if envWebRoot == "" {
		if err := os.Setenv(webRootEnvString, defaultWebRoot); err != nil {
			panic(err)
		}
		webRoot = defaultWebRoot
	} else {
		webRoot = envWebRoot
	}
}

func noCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Writer.Header().Set("Expires", "-1")

		c.Next()
	}
}

// Run starts the http server
func Run(version string) {
	server := gin.Default()

	server.Use(noCache())

	//indexFile := filepath.Join(webRoot, "index.html")
	server.Use(static.Serve("/", static.LocalFile(webRoot, true)))

	server.GET("/iceservers", iceServersHandler)

	callsGroup := server.Group("/calls")
	callsGroup.GET("", getCallsHandler)
	callsGroup.POST("", newCallHandler)

	callGroup := callsGroup.Group("/:id")
	callGroup.POST("/icecandidate", iceCandidateHandler)
	callGroup.POST("/join", joinCallHandler)
	callGroup.POST("/leave", leaveCallHandler)
	callGroup.PUT("/renegotiate", renegotiateCallHandler)

	server.GET("/websocket", websocketHandler)

	server.Run(port)
}
