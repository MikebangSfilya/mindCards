package dtoout

import "time"

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
