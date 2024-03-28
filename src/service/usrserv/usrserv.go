package usrserv

import (
	"crypto/sha512"
	"dwelt/src/dto"
	"dwelt/src/model/dao"
	"dwelt/src/model/entity"
	"dwelt/src/ws/chat"
	"encoding/hex"
	"log/slog"
)

type UserService struct {
	wsHub   *chat.Hub
	userDao *dao.UserDao
}

func NewUserService(wsHub *chat.Hub, userDao *dao.UserDao) *UserService {
	return &UserService{
		wsHub:   wsHub,
		userDao: userDao,
	}
}

func hashPassword(password string) string {
	h := sha512.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func (us *UserService) ValidateUser(username string, password string) (userId int64, valid bool, err error) {
	user, err := us.userDao.FindUserByUsernameAndPassword(username, hashPassword(password))
	if err != nil {
		slog.Error(err.Error(), "method", "ValidateUser")
		return
	}

	if user == nil {
		valid = false
		return
	}

	userId = user.ID
	valid = true

	return
}

func (us *UserService) RegisterUser(username string, password string) (userId int64, duplicate bool, err error) {
	userId, duplicate, err = us.userDao.CreateUser(username, hashPassword(password))
	return
}

func (us *UserService) SearchUsers(prefix string, limit int) (users []dto.UserResponse, err error) {
	usersEntity, err := us.userDao.SearchUsers(prefix, limit)
	if err != nil {
		return
	}

	users = make([]dto.UserResponse, len(usersEntity))
	for i, user := range usersEntity {
		users[i] = userResponseFromEntity(user)
	}

	return
}

func (us *UserService) FindDirectChat(requesterUid int64, directToUid int64) (chatId int64, err error) {
	directChat, err := us.userDao.FindDirectChat(requesterUid, directToUid)
	if err != nil {
		return
	}

	if directChat != nil {
		return directChat.ID, nil
	}

	newChat := entity.Chat{
		Name: "",
		Users: []entity.User{
			{ID: requesterUid},
			{ID: directToUid},
		},
	}

	newChatId, err := us.userDao.CreateChat(&newChat)
	if err != nil {
		return
	}

	return newChatId, nil
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
	chatEntity, err := us.userDao.FindChatById(message.Message.ChatId)
	if err != nil {
		slog.Error(err.Error(), "method", "handleMessage")
		return
	}

	// check if the user is in the chat
	userInChat := false
	username := ""
	var receiversUserIds []int64
	for _, user := range chatEntity.Users {
		receiversUserIds = append(receiversUserIds, user.ID)
		if user.ID == message.ClientId {
			userInChat = true
			username = user.Username
		}
	}

	if !userInChat {
		slog.Error("User not in chat", "method", "handleMessage")
		return
	}

	// save message
	messageEntity := entity.Message{
		Text:   message.Message.Message,
		ChatId: message.Message.ChatId,
	}

	err = us.userDao.SaveMessage(&messageEntity)
	if err != nil {
		slog.Error(err.Error(), "method", "handleMessage")
		return
	}

	serverMessage := dto.WebSocketServerMessage{
		ChatId:    message.Message.ChatId,
		UserId:    message.ClientId,
		Username:  username,
		Message:   message.Message.Message,
		CreatedAt: messageEntity.CreatedAt,
	}

	us.wsHub.SendToSelected(serverMessage, receiversUserIds)
}
