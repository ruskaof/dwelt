package chat

import (
	"dwelt/src/dto"
	"dwelt/src/metrics"
	"github.com/gorilla/websocket"
	"log"
	"log/slog"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	// maxMessageSize = 512 todo: use this
)

type IncomingClientMessage struct {
	ClientId int64
	Message  dto.WebSocketClientMessage
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userId int64
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		if err := c.conn.Close(); err != nil {
			log.Println(err)
		}
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		slog.Debug("Got pong from the client", "address", c.conn.NetConn().RemoteAddr().String())
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			slog.Debug(err.Error(), "method", "Ws read message")
			break
		}
		deserializedMessage, err := dto.DeserializeWebSocketClientMessage(message)
		if err != nil {
			slog.Debug(err.Error(), "method", "Ws deserialize message")
			continue
		}

		c.hub.Incoming <- IncomingClientMessage{
			ClientId: c.userId,
			Message:  deserializedMessage,
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := c.conn.Close(); err != nil {
			slog.Error(err.Error(), "method", "Ws connection close deffering")
		}
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				slog.Debug("Ws connection was closed", "address", c.conn.NetConn().RemoteAddr())
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				slog.Error(err.Error(), "method", "Ws next writer")
				return
			}
			w.Write(message)
			metrics.IncrementIncomingWebsocketBytes(len(message))

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
				metrics.IncrementIncomingWebsocketBytes(len(message))
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			slog.Debug("Sending ping to the client", "address", c.conn.NetConn().RemoteAddr().String())
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				slog.Error(err.Error(), "method", "Ws ping")
				return
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, userId int64, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error(err.Error(), "method", "Ws upgrade")
		return
	}
	client := &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		userId: userId,
	}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
