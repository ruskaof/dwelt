package handler

import "net/http"

func retrieveUserId(r *http.Request) int64 {
	userId, _ := r.Context().Value("userId").(int64)
	return userId
}
