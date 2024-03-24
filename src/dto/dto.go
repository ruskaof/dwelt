package dto

import "encoding/json"

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
	ChatId  int64  `json:"chatId"`
	Message string `json:"message"`
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
	ChatId int64 `json:"chatId"`
}
