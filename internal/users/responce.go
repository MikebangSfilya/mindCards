package users

import (
	"encoding/json"
	"time"
)

type UserResponce struct {
	Email  string `json:"email"`
	UserId int    `json:"user_id"`
}

type ErrResponce struct {
	Err  string
	Time time.Time
}

func NewErr(err error) ErrResponce {
	return ErrResponce{
		Err:  err.Error(),
		Time: time.Now(),
	}
}

func (e *ErrResponce) ToString() string {
	b, err := json.Marshal(e)
	if err != nil {
		return ""
	}
	return string(b)
}
