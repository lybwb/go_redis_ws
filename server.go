package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var manager = ClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
			}
		case message := <-manager.broadcast:
			fmt.Println(string(message))
			for conn := range manager.clients {
				fmt.Println("To client: ", conn.id)

				select {
				case conn.send <- message:
				default:
					fmt.Println("close connect")
					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}

func (c *Client) write() {
	defer func() {
		manager.unregister <- c
		c.socket.Close()
		fmt.Println("writer closed")
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.socket.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				manager.unregister <- c
				c.socket.Close()
				fmt.Println("write error!")
				break
			}
			fmt.Println("Write data")
		}
	}
}

func main() {
	fmt.Println("Starting application...")
	go manager.start()
	go manager.getRedisData()
	http.HandleFunc("/ws", wsPage)
	http.ListenAndServe(":12345", nil)
}

func wsPage(res http.ResponseWriter, req *http.Request) {
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		io.WriteString(res, "error connect")
		return
	}

	uid := uuid.NewV4()
	sha1 := uid.String()

	client := &Client{id: sha1, socket: conn, send: make(chan []byte)}
	manager.register <- client

	go client.write()
}

func (manager *ClientManager) getRedisData() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "test123",
		DB:       4,
	})
	redisSubscript := redisClient.Subscribe("chl")
	for {
		msg, err := redisSubscript.ReceiveMessage()
		if err != nil {
			redisClient.Close()
		}
		jsonMessage, _ := json.Marshal(&Message{Sender: "sender", Content: msg.String()})
		manager.broadcast <- jsonMessage

	}
}
