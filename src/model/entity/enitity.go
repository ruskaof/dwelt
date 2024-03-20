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
