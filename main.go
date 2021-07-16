package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocketHub is websocket connection hub
var WebSocketHub *WSHub

// WebSocketUpgrader is a global upgrader for each connectoin
var WebSocketUpgrader websocket.Upgrader

func main() {
	fmt.Println("Starting application...")

	// websocket
	WebSocketHub = SetupWebSocketHub()
	WebSocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// CheckOrigin:     WebSocketCheckOirignHandler,
	}
	go WebSocketHub.WebSocketRun()
	go WebSocketHub.WSChannelRun()

	http.HandleFunc("/ws", APIWSHandler)
	http.ListenAndServe(":12346", nil)
}
