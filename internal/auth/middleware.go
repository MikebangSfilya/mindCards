package auth

import (
	"context"
	"net/http"
	"strconv"
)

const (
	ctxKeyUser ctxKey = iota
)

type ctxKey int8

// get header userId
func AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDHead := r.Header.Get("X-User-ID")

		if userIDHead == "" {
			http.Error(w, "X-User-ID header required", http.StatusBadRequest)
			return
		}

		userID, err := strconv.Atoi(userIDHead)
		if err != nil || userID <= 0 {
			http.Error(w, "Invalid X-User-ID. Must be positive integer", http.StatusBadRequest)
			return
		}

		//TODO: replase 1 to real user
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, userID)))
	})
}
