package chat

import (
	"encoding/json"
	"forum/backend/config"
	"forum/backend/database"
	"forum/backend/structure"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hub        *Hub
	conn       *websocket.Conn // websocket connections
	send       chan []byte
	userID     int
	typing     bool
	typingLock chan bool
	isReceiver bool
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error ReadPump: %v", err)
			}
			break
		}

		log.Printf("Message recu depuis le ien-cli: %s", string(message))

		var msg structure.Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			break
		}

		msg.Sender_id = c.userID

		if msg.Msg_type == "msg" {
			msg.Date = time.Now().Format("01-02-2006 15:04:05")

			err = database.NewMessage(config.Path, msg)
			if err != nil {
				log.Printf("Error storing new message: %v", err)
				break
			}
		}
		sendMsg, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error marshalling message: %v", err)
			break
		}

		c.hub.broadcast <- sendMsg
	}

	c.typingLock <- false
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case isTyping := <-c.typingLock:
			c.typing = isTyping

			if c.isReceiver {
				typingStatus := struct {
					ReceiverID int    `json:"receiver_id"`
					IsTyping   bool   `json:"is_typing"`
					MsgType    string `json:"msg_type"`
				}{
					ReceiverID: c.userID,
					IsTyping:   c.typing,
					MsgType:    "typing",
				}
				sendTypingStatus, err := json.Marshal(typingStatus)
				if err != nil {
					log.Println("Error marshalling typing status:", err)
					continue
				}
				c.hub.broadcast <- sendTypingStatus
			}
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		return
	}

	foundVal := cookie.Value

	curr, err := database.CurrentUser(config.Path, foundVal)
	if err != nil {
		return
	}

	client := &Client{
		hub:        hub,
		conn:       conn,
		send:       make(chan []byte, 256),
		userID:     curr.ID,
		typing:     false,
		typingLock: make(chan bool),
		isReceiver: false,
	}

	log.Println("Client isReceiver:", client.isReceiver)
	client.hub.register <- client

	go client.writePump()
	go client.ReadPump()

}
