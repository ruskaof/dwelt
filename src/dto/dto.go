package dto

import (
	"encoding/json"
	"time"
)

type UserInfo struct {
	UserId int64 `json:"id"`
}

type UserResponse struct {
	UserId   int64  `json:"userId"`
	Username string `json:"username"`
}

type WebSocketClientMessage struct {
	ChatId  int64  `json:"chatId"`
	Message string `json:"message"`
}

type WebSocketServerMessage struct {
	ChatId    int64     `json:"chatId"`
	UserId    int64     `json:"userId"`
	Username  string    `json:"username"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

func SerializeWebSocketServerMessage(message WebSocketServerMessage) []byte {
	res, _ := json.Marshal(message)
	return res
}

func DeserializeWebSocketServerMessage(data []byte) (message WebSocketServerMessage, err error) {
	err = json.Unmarshal(data, &message)
	return
}

func SerializeWebSocketClientMessage(message WebSocketClientMessage) []byte {
	res, _ := json.Marshal(message)
	return res
}

func DeserializeWebSocketClientMessage(data []byte) (message WebSocketClientMessage, err error) {
	err = json.Unmarshal(data, &message)
	return
}

type OpenDirectChatResponse struct {
	ChatId       int64                    `json:"chatId"`
	LastMessages []WebSocketServerMessage `json:"lastMessagesSorted"`
}
