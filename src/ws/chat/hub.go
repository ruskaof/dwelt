package chat

import "dwelt/src/dto"

type Hub struct {
	clients    map[*Client]bool
	Incoming   chan IncomingClientMessage
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Incoming:   make(chan IncomingClientMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		}
	}
}

func (h *Hub) SendToSelected(message dto.WebSocketServerMessage, ids []int64) {
	for client := range h.clients {
		for _, id := range ids {
			if client.userId == id {
				client.send <- dto.SerializeWebSocketServerMessage(message)
			}
		}
	}
}
