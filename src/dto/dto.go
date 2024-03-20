package dto

type UserInfo struct {
	UserId int64 `json:"id"`
}

type UserResponse struct {
	UserId   int64  `json:"userId"`
	Username string `json:"username"`
}
