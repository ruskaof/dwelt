package entity

type UserEntity struct {
	ID       int32  `gorm:"column:id"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}
