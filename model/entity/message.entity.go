package entity

import (
	"github.com/gofiber/contrib/websocket"
)

type Message struct {
	Name    string
	Message string
}

type Hub struct {
	clients map[*websocket.Conn]bool
	clientRegisterChannel chan *websocket.Conn
	clientRemovalChannel chan *websocket.Conn
	broadcastMessage chan Message
}

func (h *Hub) Run() {
	for {
		select {
		case conn:= <- h.clientRegisterChannel:
			h.clients[conn] = true
		case conn := <- h.clientRemovalChannel:
			delete(h.clients, conn)
		case msg := <- h.broadcastMessage:
			for conn := range h.clients {
				_ = conn.WriteJSON(msg)
			}
		}

	}
}

func NewHub() *Hub {
    return &Hub{
        clients:               make(map[*websocket.Conn]bool),
        clientRegisterChannel: make(chan *websocket.Conn),
        clientRemovalChannel:  make(chan *websocket.Conn),
        broadcastMessage:     make(chan Message),
    }
}

func SendMessage(h *Hub) func (*websocket.Conn) {
	return func(conn *websocket.Conn) {
		defer func ()  {
			h.clientRemovalChannel <- conn
			_ = conn.Close()
		}()

		name := conn.Query("name", "")
		h.clientRegisterChannel <- conn

		for {
			messageType, text, err := conn.ReadMessage()
			if err != nil {
				return
			}

			if messageType == websocket.TextMessage {
				h.broadcastMessage <- Message{
					Name: name,
					Message: string(text),
				}
			}
		}
	}
}