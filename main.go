package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/gorilla/websocket"
)

// WebSocketHub is websocket connection hub
var WebSocketHub *WSHub
var TUid int

// WebSocketUpgrader is a global upgrader for each connectoin
var WebSocketUpgrader websocket.Upgrader

func main() {
	fmt.Println("Starting application...")

	port := os.Args[1]

	uid, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println(err)
	}

	TUid = uid
	fmt.Println("port", port)
	fmt.Println("test uid:", TUid)

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

	// http.ListenAndServe(":12346", nil)
	router.Run(":" + port)
}
