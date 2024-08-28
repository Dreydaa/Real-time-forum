package chat

import (
	"encoding/json"
	"forum/backend/structure"
)

type Hub struct {
	clients      map[int]*Client
	broadcast    chan []byte
	register     chan *Client
	unregister   chan *Client
	typing       map[int]bool
	typing2      map[string]bool
	typingStatus map[int]int
}

func NewHub() *Hub {
	return &Hub{
		broadcast:    make(chan []byte),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		clients:      make(map[int]*Client),
		typing:       make(map[int]bool),
		typing2:      make(map[string]bool),
		typingStatus: make(map[int]int),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.userID] = client

			uids := make([]int, 0, len(h.clients))
			for id := range h.clients {
				uids = append(uids, id)
			}
			msg := structure.OnlineUsers{
				UserIds:  uids,
				Msg_type: "online",
			}
			sendMsg, err := json.Marshal(msg)
			if err != nil {
				panic(err)
			}

			for _, c := range h.clients {
				select {
				case c.send <- sendMsg:
				default:
					close(c.send)
					delete(h.clients, c.userID)
				}
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)

				uids := make([]int, 0, len(h.clients))
				for id := range h.clients {
					uids = append(uids, id)
				}
				msg := structure.OnlineUsers{
					UserIds:  uids,
					Msg_type: "offline",
				}
				sendMsg, err := json.Marshal(msg)
				if err != nil {
					panic(err)
				}

				for _, c := range h.clients {
					select {
					case c.send <- sendMsg:
					default:
						close(c.send)
						delete(h.clients, c.userID)
					}
				}
				close(client.send)
			}
		case message := <-h.broadcast:
			var msg structure.Message
			if err := json.Unmarshal(message, &msg); err != nil {
				panic(err)
			}

			sendMsg, err := json.Marshal(msg)
			if err != nil {
				panic(err)
			}

			if msg.Msg_type == "msg" {
				for _, client := range h.clients {
					if client.userID == msg.Receiver_id {
						select {
						case client.send <- sendMsg:
						default:
							close(client.send)
							delete(h.clients, client.userID)
						}
					}
				}
			} else {
				for _, client := range h.clients {
					if client.userID != msg.Sender_id {
						select {
						case client.send <- sendMsg:
						default:
							close(client.send)
							delete(h.clients, client.userID)
						}
					}
				}
			}
		}
	}
}
