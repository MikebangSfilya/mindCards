package users

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type Save interface {
	SaveUser(ctx context.Context, user *User) error
}

func SaveUser(save Save) http.HandlerFunc {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		var uReq Request
		if err := decoder(r, &uReq); err != nil {
			handleError(w, err, err.Error(), http.StatusBadRequest)
			return
		}
		u, err := NewUser(uReq.Email, uReq.Password)

		if err != nil {
			handleError(w, err, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := save.SaveUser(ctx, u); err != nil {
			handleError(w, err, err.Error(), http.StatusBadRequest)
			return
		}

		result := UserResponce{
			Email:  u.Email,
			UserId: u.UserId,
		}

		if err := encoder(w, result); err != nil {
			handleError(w, err, err.Error(), http.StatusBadRequest)
			return
		}

	}
}

// help func for decode json
func decoder(r *http.Request, dto any) error {
	return json.NewDecoder(r.Body).Decode(dto)
}

// help func for encode json
func encoder(w http.ResponseWriter, resp any) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resp)
}

// help func responce error DTO
func handleError(w http.ResponseWriter, err error, msg string, code int) {
	slog.Error(msg, "err", err, "package", "handlers")
	errDTO := NewErr(err)
	http.Error(w, errDTO.ToString(), code)
}
