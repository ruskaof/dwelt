package usrserv

import (
	"dwelt/src/model/entity"
)

type UserResponse struct {
	UserId   int64  `json:"userId"`
	Username string `json:"username"`
}

func userResponseFromEntity(user entity.User) UserResponse {
	return UserResponse{
		UserId:   user.ID,
		Username: user.Username,
	}
}
