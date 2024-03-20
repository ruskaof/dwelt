package usrserv

import (
	"dwelt/src/dto"
	"dwelt/src/model/entity"
)

func userResponseFromEntity(user entity.User) dto.UserResponse {
	return dto.UserResponse{
		UserId:   user.ID,
		Username: user.Username,
	}
}
