package middleware

import (
	"log"
	"net/http"
)

func CheckAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(UserID(UserIDContextKey))
		log.Println(userId)
		if userId == nil {
			http.Error(w, "User not logged in", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
