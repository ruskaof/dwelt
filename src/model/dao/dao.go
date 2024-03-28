package dao

import (
	"dwelt/src/config"
	"dwelt/src/model/entity"
	"dwelt/src/utils"
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.DbCfg.Host,
		config.DbCfg.User,
		config.DbCfg.Password,
		config.DbCfg.DbName,
		config.DbCfg.Port,
	)

	return utils.Must(
		gorm.Open(postgres.Open(dsn), &gorm.Config{
			TranslateError: true,
			// log every SQL command
			Logger: logger.Default.LogMode(logger.Info),
		}),
	)
}

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (ud *UserDao) FindUserByUsernameAndPassword(username, password string) (*entity.User, error) {
	users := make([]entity.User, 0)
	err := ud.db.Where("username = ? AND password = ?", username, password).Find(&users).Error
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, nil
	}

	return &users[0], nil
}

func (ud *UserDao) FindUserById(userId int64) (*entity.User, error) {
	users := make([]entity.User, 0)
	err := ud.db.Where("id = ?", userId).Find(&users).Error
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, nil
	}

	return &users[0], nil
}

// CreateUser creates a new user in the database
func (ud *UserDao) CreateUser(username, password string) (userId int64, duplicate bool, err error) {
	user := entity.User{
		Username: username,
		Password: password,
	}

	err = ud.db.Create(&user).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return 0, true, nil
	}

	return user.ID, false, err
}

// SearchUsers searches for users with a username that starts with the given prefix
func (ud *UserDao) SearchUsers(prefix string, limit int) (users []entity.User, err error) {
	err = ud.db.Where("username LIKE ?", prefix+"%").Limit(limit).Find(&users).Error
	return
}

func (ud *UserDao) FindDirectChat(userId1, userId2 int64) (*entity.Chat, error) {
	var chatEntity entity.Chat
	err := ud.db.
		Joins("JOIN users_chats uc1 ON chats.id = uc1.chat_id").
		Joins("JOIN users_chats uc2 ON uc1.chat_id = uc2.chat_id").
		Where("uc1.user_id = ?", userId1).
		Where("uc2.user_id = ?", userId2).
		First(&chatEntity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &chatEntity, nil
}

// CreateChat creates a new chat in the database
func (ud *UserDao) CreateChat(chat *entity.Chat) (chatId int64, err error) {
	err = ud.db.Create(chat).Error
	return chat.ID, err
}

func (ud *UserDao) FindChatById(chatId int64) (*entity.Chat, error) {
	chat := &entity.Chat{
		ID: chatId,
	}
	err := ud.db.Preload("Users").First(chat).Error
	return chat, err
}

// SaveMessage saves a message in the database
func (ud *UserDao) SaveMessage(message *entity.Message) error {
	return ud.db.Create(message).Error
}

func (ud *UserDao) FindLastMessagesByChat(chatId int64, limit int32) (messages []entity.Message, err error) {
	err = ud.db.
		Preload("User").
		Where("chat_id = ?", chatId).Order("created_at desc").Limit(int(limit)).Find(&messages).Error
	return
}
