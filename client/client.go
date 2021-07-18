package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

func main() {

	port := os.Args[1]
	fmt.Println("port", port)

	uid := os.Args[2]

	var addr = flag.String("addr", "localhost:"+port, "http service address")

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	var dialer *websocket.Dialer

	// conn, _, err := dialer.Dial(u.String(), nil)

	conn, _, err := dialer.Dial(u.String(), http.Header{"x-uid": []string{uid}})
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}

		fmt.Printf("received: %s\n", message)
	}
}
