package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/gorilla/websocket"
)

// WebSocketHub is websocket connection hub
var WebSocketHub *WSHub

// WebSocketUpgrader is a global upgrader for each connectoin
var WebSocketUpgrader websocket.Upgrader

func main() {

	port := os.Args[1]
	fmt.Println("Starting application..port.", port)

	// websocket
	WebSocketHub = SetupWebSocketHub()
	WebSocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     WebSocketCheckOirignHandler,
	}
	go WebSocketHub.WebSocketRun()
	go WebSocketHub.WSChannelRun()

	router := gin.Default()

	// http.HandleFunc("/ws", APIWSHandler)
	router.GET("/ws", APIWSHandler)
	router.GET("/pub", APIWSPublishHandler)

	// http.ListenAndServe(":12346", nil)
	router.Run(":" + port)
}
