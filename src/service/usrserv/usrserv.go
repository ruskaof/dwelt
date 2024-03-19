package usrserv

import (
	"crypto/sha512"
	"dwelt/src/model/dao"
	"dwelt/src/model/entity"
	"encoding/hex"
	"errors"
	"gorm.io/gorm"
	"log/slog"
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

func SearchUsers(substring string, limit int) (users []UserResponse, err error) {
	var usersEntity []entity.User
	err = dao.Db.Where("username LIKE ?", "%"+substring+"%").Limit(limit).Find(&usersEntity).Error
	if err != nil {
		slog.Error(err.Error(), "method", "SearchUsers")
	}
	
	users = make([]UserResponse, len(usersEntity))
	for i, user := range usersEntity {
		users[i] = userResponseFromEntity(user)
	}

	return
}