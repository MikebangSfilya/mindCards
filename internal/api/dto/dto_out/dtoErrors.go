package dtoout

import (
	"encoding/json"
	"log"
	"time"
)

type ErrDto struct {
	Err  string
	Time time.Time
}

func NewErr(err error) ErrDto {
	return ErrDto{
		Err:  err.Error(),
		Time: time.Now(),
	}
}

func (e *ErrDto) ToString() string {
	b, err := json.Marshal(e)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(b)
}
