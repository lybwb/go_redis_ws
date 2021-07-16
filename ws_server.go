package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ecopia-china/hdmap-platform/backend/commons"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	// maximum connection number for one user.
	maxUserConns = 10
)

// Websocket message transfering format
type WSMessage struct {
	Message    NotifyUserMessage `json:"message"`
	NewMessage bool              `json:"new_message"`
	UserNotifyUnread

	// type Message struct {
	// 	Sender    string `json:"sender,omitempty"`
	// 	Recipient string `json:"recipient,omitempty"`
	// 	Content   string `json:"content,omitempty"`
	// }
}

// Websocket message connection format
type WSMessageConn struct {
	Uid         int
	MessageJSON []byte
}

type WSHub struct {
	// Registered clients.
	clients map[*WSClient]bool

	// map client user uid to user's related WSClient
	clientUidMap map[int]map[int]*WSClient

	// Inbound messages from the clients.
	broadcast chan []byte

	// messages send from server to connection
	toMessage chan WSMessageConn

	// Register requests from the clients.
	register chan *WSClient

	// Unregister requests from clients.
	unregister chan *WSClient

	// type ClientManager struct {
	// clients    map[*Client]bool
	// broadcast  chan []byte
	// register   chan *Client
	// unregister chan *Client
	// }
}

type WSClient struct {
	hub *WSHub

	// The websocket connection.
	conn *websocket.Conn

	// connection id within the same user(uid)
	connID int

	// Buffered channel of outbound messages.
	send chan []byte

	// user's uid
	uid int

	// type Client struct {
	// 	id     string
	// 	socket *websocket.Conn
	// 	send   chan []byte
}

func SetupWebSocketHub() *WSHub {
	return &WSHub{
		clients:      make(map[*WSClient]bool),
		clientUidMap: make(map[int]map[int]*WSClient),
		broadcast:    make(chan []byte),
		toMessage:    make(chan WSMessageConn),
		register:     make(chan *WSClient),
		unregister:   make(chan *WSClient),
	}

	// var manager = ClientManager{
	// 	broadcast:  make(chan []byte),
	// 	register:   make(chan *Client),
	// 	unregister: make(chan *Client),
	// 	clients:    make(map[*Client]bool),
	// }
}

func WebSocketCheckOirignHandler(req *http.Request) bool {
	return false
}

func (h *WSHub) WebSocketRun() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			if _, ok := h.clientUidMap[client.uid]; !ok {
				h.clientUidMap[client.uid] = make(map[int]*WSClient)
			}
			h.clientUidMap[client.uid][client.connID] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.clientUidMap[client.uid], client.connID)
				if len(h.clientUidMap[client.uid]) <= 0 {
					delete(h.clientUidMap, client.uid)
				}
				close(client.send)
			}
		case message := <-h.toMessage:
			clientConnIDs, ok := h.clientUidMap[message.Uid]
			if ok {
				for _, client := range clientConnIDs {
					if h.clients[client] {
						client.send <- message.MessageJSON
					}
				}
			}
		}
	}
}

func (h *WSHub) WSChannelRun() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "test123", // no password set
		DB:       4,         // use default DB
	})
	// redisSubscript := redisClient.Subscribe("mychannel1")
	redisSubscript := redisClient.PSubscribe("*")

	for {

		for msg := range redisSubscript.Channel() {
			fmt.Printf("channel=%s message=%s\n", msg.Channel, msg.Payload)
			jsonMessage, _ := json.Marshal(&Message{Sender: "hi", Content: msg.String()})
			manager.broadcast <- jsonMessage
		}

	}
}

func (h *WSHub) send(message []byte, ignore *Client) {
	for conn := range h.clients {
		if conn != ignore {
			conn.send <- message
		}
	}
}

func (c *WSClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(commons.GetNowUTCTime().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(commons.GetNowUTCTime().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func serveWebSocket(hub *WSHub, w http.ResponseWriter, r *http.Request, user User) error {
	// func APIWSHandler(res http.ResponseWriter, req *http.Request) {
	conn, err := WebSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	connID := GetUserClientConnID(hub, user.UserID)
	// client := &Client{id: sha1, socket: conn, send: make(chan []byte)}
	client := &WSClient{
		hub:    hub,
		conn:   conn,
		connID: connID,
		send:   make(chan []byte, 256),
		uid:    user.UserID,
	}
	client.hub.register <- client

	// go client.write()
	go client.writePump()

	_, _, _, err = SendWSMessageToConn(NotifyUserMessage{}, user.UserID, 0, 0, 0)
	if err != nil {
		LogError("serveWebSocket: SendWSMessageToConn fail", zap.Error(err))
		return err
	}

	return nil

}

func APIWSHandler(c *gin.Context) {

	// user, err := isTokenAuthorized(accessToken)

	err = serveWebSocket(WebSocketHub, c.Writer, c.Request, user)
	if err != nil {
		LogError("ServeWebSocketHandler: call serveWebSocket fail", zap.Error(err))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

}
