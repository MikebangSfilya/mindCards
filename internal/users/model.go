package users

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var errValidate = fmt.Errorf("zero len email or password")

type User struct {
	Email             string `json:"email"`
	EncryptedPassword string `json:"-"`
	UserId            int    `json:"user_id"`
}

func NewUser(email, pass string) (*User, error) {
	if !simpleValidate(email, pass) {
		return nil, errValidate
	}

	u := &User{
		Email: email,
	}

	if len(pass) > 0 {
		u.EncryptedPassword = encryptPass(pass)
	}

	return u, nil
}

func simpleValidate(email, pass string) bool {
	if len(email) == 0 {
		return false
	}
	if len(pass) == 0 {
		return false
	}
	return true
}

func encryptPass(pass string) string {
	passCr, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		return ""
	}
	return string(passCr)
}
