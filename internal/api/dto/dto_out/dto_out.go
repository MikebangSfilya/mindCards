package dtoout

import "time"

type DTOdel_out struct {
	Status  string    `json:"Status"`
	Object  string    `json:"Object"`
	TimeDel time.Time `json:"TimeDel"`
}

func NewDelDTO(obj string) DTOdel_out {
	return DTOdel_out{
		Status:  "deleted",
		Object:  obj,
		TimeDel: time.Now(),
	}
}
