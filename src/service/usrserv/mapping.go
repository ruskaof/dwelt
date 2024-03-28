package usrserv

import (
	"dwelt/src/dto"
	"dwelt/src/model/entity"
)

func userResponseFromEntity(user entity.User) dto.UserResponse {
	return dto.UserResponse{
		UserId:   user.ID,
		Username: user.Username,
	}
}

func messageEntityToWebSocketServerMessage(message entity.Message) dto.WebSocketServerMessage {
	return dto.WebSocketServerMessage{
		ChatId:    message.ChatId,
		UserId:    message.UserId,
		Username:  message.User.Username,
		Message:   message.Text,
		CreatedAt: message.CreatedAt,
	}
}

func mapMessagesToWebSocketServerMessages(messages []entity.Message) []dto.WebSocketServerMessage {
	var res []dto.WebSocketServerMessage
	for _, message := range messages {
		res = append(res, messageEntityToWebSocketServerMessage(message))
	}
	return res
}
