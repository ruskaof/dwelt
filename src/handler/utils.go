package handler

import "net/http"

type middleware func(http.HandlerFunc) http.HandlerFunc

func makeHandler(h http.HandlerFunc, middlewares ...middleware) http.HandlerFunc {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}

func retrieveUserId(r *http.Request) int64 {
	userId, _ := r.Context().Value("userId").(int64)
	return userId
}
