package entity

import "time"

type User struct {
	ID       int64  `gorm:"column:id"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
	Chats    []Chat `gorm:"many2many:users_chats;"`
}

type Chat struct {
	ID        int64     `gorm:"column:id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	// Name is empty for a direct chat
	Name  string `gorm:"column:name"`
	Users []User `gorm:"many2many:users_chats;"`
}

type Message struct {
	ID        int64     `gorm:"column:id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	Text      string    `gorm:"column:text"`
	ChatId    int64     `gorm:"column:chat_id"`
	UserId    int64     `gorm:"column:user_id"`
	User      User      `gorm:"foreignkey:UserId"`
}
