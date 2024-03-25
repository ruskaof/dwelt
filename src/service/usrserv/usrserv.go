package usrserv

import (
	"crypto/sha512"
	"dwelt/src/dto"
	"dwelt/src/model/entity"
	"dwelt/src/ws/chat"
	"encoding/hex"
	"errors"
	. "github.com/samber/lo"
	"gorm.io/gorm"
	"log/slog"
)

type UserService struct {
	wsHub *chat.Hub
	db    *gorm.DB
}

func NewUserService(wsHub *chat.Hub, db *gorm.DB) *UserService {
	return &UserService{
		wsHub: wsHub,
		db:    db,
	}
}

func hashPassword(password string) string {
	h := sha512.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func (us *UserService) ValidateUser(username string, password string) (userId int64, valid bool, err error) {
	var user entity.User

	err = us.db.Where("username = ? AND password = ?", username, hashPassword(password)).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Debug("User not found", "username", username, "password", password)
		err = nil
		return
	}
	if err != nil {
		slog.Error(err.Error(), "method", "ValidateUser")
		return
	}

	userId = user.ID
	valid = true

	return
}

func (us *UserService) RegisterUser(username string, password string) (userId int64, duplicate bool, err error) {
	user := entity.User{
		Username: username,
		Password: hashPassword(password),
	}

	err = us.db.Create(&user).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		duplicate = true
		err = nil
		return
	}
	if err != nil {
		slog.Error(err.Error(), "method", "RegisterUser")
		return
	}

	userId = user.ID
	return
}

func (us *UserService) SearchUsers(prefix string, limit int) (users []dto.UserResponse, err error) {
	var usersEntity []entity.User
	err = us.db.Where("username LIKE ?", prefix+"%").Limit(limit).Find(&usersEntity).Error
	if err != nil {
		slog.Error(err.Error(), "method", "SearchUsers")
	}

	users = make([]dto.UserResponse, len(usersEntity))
	for i, user := range usersEntity {
		users[i] = userResponseFromEntity(user)
	}

	return
}

func (us *UserService) FindDirectChat(requesterUid int64, directToUid int64) (chatId int64, badUsers bool, err error) {
	// check if both users exist
	var count int64
	err = us.db.Model(&entity.User{}).Where("id IN (?)", []int64{requesterUid, directToUid}).Count(&count).Error
	if err != nil {
		slog.Error(err.Error(), "method", "FindDirectChat")
		return
	}
	if count != 2 {
		badUsers = true
		return
	}

	var chatEntity entity.Chat
	err = us.db.
		Joins("JOIN users_chats uc1 ON chats.id = uc1.chat_id").
		Joins("JOIN users_chats uc2 ON uc1.chat_id = uc2.chat_id").
		Where("uc1.user_id = ?", requesterUid).
		Where("uc2.user_id = ?", directToUid).
		First(&chatEntity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	if err != nil {
		slog.Error(err.Error(), "method", "CreateDirectChat")
		return
	}

	if chatEntity.ID != 0 {
		chatId = chatEntity.ID
		return
	}

	// create chat with associated users but don't create the users
	chatEntity = entity.Chat{
		Users: []entity.User{
			{ID: requesterUid},
			{ID: directToUid},
		},
	}
	err = us.db.Create(&chatEntity).Error
	if err != nil {
		slog.Error(err.Error(), "method", "CreateDirectChat")
		return
	}

	chatId = chatEntity.ID
	return
}

func (us *UserService) StartHandlingMessages() {
	go func() {
		for {
			select {
			case message := <-us.wsHub.Incoming:
				slog.Debug("Handling message", "message", message)
				us.handleMessage(message)
			}
		}
	}()
}

func (us *UserService) handleMessage(message chat.IncomingClientMessage) {
	// find chat
	chatEntity := entity.Chat{
		ID: message.Message.ChatId,
	}

	err := us.db.Model(&entity.Chat{}).Preload("Users").First(&chatEntity, message.Message.ChatId).Error
	if err != nil {
		slog.Error(err.Error(), "method", "HandleMessage")
		return
	}

	// check if user is in chat
	inChat := false
	for _, user := range chatEntity.Users {
		if user.ID == message.ClientId {
			inChat = true
			break
		}
	}

	if !inChat {
		slog.Error("User is not in chat", "userId", message.ClientId, "chatId", message.Message.ChatId)
		return
	}

	serverMessage := dto.WebSocketServerMessage{
		ChatId:  message.Message.ChatId,
		Message: message.Message.Message,
	}

	us.wsHub.SendToSelected(
		serverMessage,
		Map(chatEntity.Users,
			func(user entity.User, _ int) int64 {
				return user.ID
			}),
	)
}
