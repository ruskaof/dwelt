package handler

import "net/http"

type middleware func(http.HandlerFunc) http.HandlerFunc

func retrieveUserId(r *http.Request) int64 {
	userId, _ := r.Context().Value("userId").(int64)
	return userId
}
