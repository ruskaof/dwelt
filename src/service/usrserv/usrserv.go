package usrserv

import (
	"crypto/sha512"
	"dwelt/src/dto"
	"dwelt/src/model/dao"
	"dwelt/src/model/entity"
	"encoding/hex"
	"errors"
	"log/slog"

	"gorm.io/gorm"
)

func hashPassword(password string) string {
	h := sha512.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func ValidateUser(username string, password string) (userId int64, valid bool, err error) {
	var user entity.User

	err = dao.Db.Where("username = ? AND password = ?", username, hashPassword(password)).First(&user).Error
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

func RegisterUser(username string, password string) (userId int64, duplicate bool, err error) {
	user := entity.User{
		Username: username,
		Password: hashPassword(password),
	}

	err = dao.Db.Create(&user).Error
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

func SearchUsers(prefix string, limit int) (users []dto.UserResponse, err error) {
	var usersEntity []entity.User
	err = dao.Db.Where("username LIKE ?", prefix+"%").Limit(limit).Find(&usersEntity).Error
	if err != nil {
		slog.Error(err.Error(), "method", "SearchUsers")
	}

	users = make([]dto.UserResponse, len(usersEntity))
	for i, user := range usersEntity {
		users[i] = userResponseFromEntity(user)
	}

	return
}

func FindDirectChat(requesterUid int64, directToUid int64) (chatId int64, badUsers bool, err error) {
	// check if both users exist
	var count int64
	err = dao.Db.Model(&entity.User{}).Where("id IN (?)", []int64{requesterUid, directToUid}).Count(&count).Error
	if err != nil {
		slog.Error(err.Error(), "method", "FindDirectChat")
		return
	}
	if count != 2 {
		badUsers = true
		return
	}

	var chat entity.Chat
	err = dao.Db.
		Joins("JOIN users_chats uc1 ON chats.id = uc1.chat_id").
		Joins("JOIN users_chats uc2 ON uc1.chat_id = uc2.chat_id").
		Where("uc1.user_id = ?", requesterUid).
		Where("uc2.user_id = ?", directToUid).
		First(&chat).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	if err != nil {
		slog.Error(err.Error(), "method", "CreateDirectChat")
		return
	}

	if chat.ID != 0 {
		chatId = chat.ID
		return
	}

	// create chat with associated users but don't create the users
	chat = entity.Chat{
		Users: []entity.User{
			{ID: requesterUid},
			{ID: directToUid},
		},
	}
	err = dao.Db.Create(&chat).Error
	if err != nil {
		slog.Error(err.Error(), "method", "CreateDirectChat")
		return
	}

	chatId = chat.ID
	return
}
