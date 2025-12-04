package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

const (
	CtxKeyUser CtxKey = iota
)

var (
	errUserNotFound = fmt.Errorf("user_id not found in context")
)

type CtxKey int8

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

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxKeyUser, userID)))
	})
}

func GetUserID(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(CtxKeyUser).(int)
	if !ok {
		return 0, errUserNotFound
	}
	return userID, nil
}
