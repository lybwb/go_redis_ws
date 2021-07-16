package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
}

// Websocket message connection format
type WSMessageConn struct {
	Uid         int
	MessageJSON []byte
}

// Hub maintains the set of active clients and messages to the clients.
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
}

// Client is a middleman between the websocket connection and the hub.
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
}

func WebSocketCheckOirignHandler(req *http.Request) bool {
	requestOrigin := req.Header.Get("Origin")
	for _, corsDomain := range configFile.CORSDOMAINS {
		if corsDomain == requestOrigin {
			return true
		}
	}
	LogWarn("WebSocket Origin Deny.", zap.String("origin", requestOrigin),
		zap.String("allow-cors", fmtString(configFile.CORSDOMAINS)))

	return false
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
				LogDebug("send message", zap.Int("receive_uid", message.Uid),
					zap.String("receive_conn", fmtString(clientConnIDs)),
					zap.String("message", string(message.MessageJSON)))
				for _, client := range clientConnIDs {
					if h.clients[client] {
						client.send <- message.MessageJSON
					}
				}
			}
		case message := <-h.broadcast:
			// broadcast is not used
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
					delete(h.clientUidMap, client.uid)
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
			wsMessage := WSMessageConn{
				Uid:         msg.Channel,
				MessageJSON: messageInByte,
			}
			WebSocketHub.toMessage <- wsMessage
		}

	}

}

func SendWSMessageToConn(
	message NotifyUserMessage, uid, unreadNoticePlus, unreadMessagePlus int, unhandledContentRequestPlus int,
) (int, int, int, error) {
	unreadNotice, unreadMessage, err := GetUnreadCountInNotifyType(uid)
	if err != nil {
		LogError("get unread count in notify type error", zap.Error(err))
		return 0, 0, 0, err
	}

	unhandledContentRequestCount, err := getUnhandledRequestCount(uid)
	if err != nil {
		LogError("get unhandled count in content_request error", zap.Error(err))
		return 0, 0, 0, err
	}

	unreadNotice += unreadNoticePlus
	unreadMessage += unreadMessagePlus
	unhandledContentRequestCount += unhandledContentRequestPlus

	wsM := WSMessage{
		Message:    message,
		NewMessage: message.MessageID != "",
	}
	wsM.UnreadNoticeCount = unreadNotice
	wsM.UnreadMessageCount = unreadMessage
	wsM.UnhandledContentRequestCount = &unhandledContentRequestCount

	messageInByte, err := json.Marshal(wsM)
	if err != nil {
		LogError("message marshal error", zap.Error(err))
		return 0, 0, 0, err
	}

	wsMessage := WSMessageConn{
		Uid:         uid,
		MessageJSON: messageInByte,
	}

	WebSocketHub.toMessage <- wsMessage

	return unreadNotice, unreadMessage, unhandledContentRequestCount, nil
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *WSClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(commons.GetNowUTCTime().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(commons.GetNowUTCTime().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure,
			) {
				LogError("unexpected Websocket close.", zap.String("message", fmtString(message)), zap.Error(err))
			}

			break
		}
		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		//c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
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

//func deleteUserClient(hub *WSHub, uid int) {
//	var client *WSClient
//	var ok bool
//	if client, ok = hub.clientUidMap[uid]; ok {
//		//fmt.Println("deleting previous uid client")
//		close(client.send)
//		delete(hub.clients, client)
//		delete(hub.clientUidMap, uid)
//	}
//}

func checkUserClientConnExeced(hub *WSHub, uid int) bool {
	if _, ok := hub.clientUidMap[uid]; !ok {
		return false
	}
	return len(hub.clientUidMap[uid]) >= maxUserConns
}

// todo optimize
func GetUserClientConnID(hub *WSHub, uid int) int {
	maxID := 0
	for connID, _ := range hub.clientUidMap[uid] {
		if connID > maxID {
			maxID = connID
		}
	}

	return maxID + 1
}

// serveWs handles websocket requests from the peer.
func serveWebSocket(hub *WSHub, w http.ResponseWriter, r *http.Request, user User) error {
	//fmt.Println("in serve WebSocket, request header connection:", r.Header.Get("Connection"))
	LogInfo("serveWebSocket: upgrade")
	conn, err := WebSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		LogError("serveWebSocket: WebSocketUpgrader.Upgrade fail", zap.Error(err))
		return err
	}

	connID := GetUserClientConnID(hub, user.UserID)
	client := &WSClient{
		hub:    hub,
		conn:   conn,
		connID: connID,
		send:   make(chan []byte, 256),
		uid:    user.UserID,
	}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()

	_, _, _, err = SendWSMessageToConn(NotifyUserMessage{}, user.UserID, 0, 0, 0)
	if err != nil {
		LogError("serveWebSocket: SendWSMessageToConn fail", zap.Error(err))
		return err
	}

	return nil
}

func APIWSHandler(c *gin.Context) {
	LogInfo("ServeWebSocketHandler: websocket client connected")
	authCookie, err := c.Request.Cookie("X-Authorization")
	if err != nil {
		LogError("ServeWebSocketHandler: get cookie fail", zap.Error(err), zap.String("authCookie", authCookie.Value))
		c.JSON(http.StatusForbidden, "Unauthorized")
		return
	}
	accessToken := strings.Replace(authCookie.Value, "Bearer ", "", 1)
	if len(accessToken) == 0 {
		LogError("ServeWebSocketHandler: accessToken null", zap.String("authCookie", authCookie.Value))
		c.JSON(http.StatusForbidden, "Unauthorized")
		return
	}
	user, err := isTokenAuthorized(accessToken)
	if err != nil {
		LogError("ServeWebSocketHandler: token auth fail", zap.Error(err), zap.String("authCookie", authCookie.Value))
		c.JSON(http.StatusForbidden, "Unauthorized")
		return
	}
	LogInfo("ServeWebSocketHandler: auth success", zap.Int("uid", user.UserID), zap.String("email", user.Email))

	if checkUserClientConnExeced(WebSocketHub, user.UserID) {
		LogError("ServeWebSocketHandler: max user connection limit reached", zap.Int("uid", user.UserID))
		c.JSON(http.StatusForbidden, "Maximum user connection limit has been reached.")
		return
	}

	err = serveWebSocket(WebSocketHub, c.Writer, c.Request, user)
	if err != nil {
		LogError("ServeWebSocketHandler: call serveWebSocket fail", zap.Error(err))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
}
