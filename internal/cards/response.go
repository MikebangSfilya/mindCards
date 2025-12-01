package cards

import (
	"encoding/json"
	"time"
)

type DTOdel_out struct {
	Status  string    `json:"Status"`
	Object  string    `json:"Object"`
	TimeDel time.Time `json:"TimeDel"`
}

type MDAddedDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Tag         string `json:"tag"`
}

func NewDelDTO(obj string) DTOdel_out {
	return DTOdel_out{
		Status:  "deleted",
		Object:  obj,
		TimeDel: time.Now(),
	}
}

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
		return ""
	}
	return string(b)
}
