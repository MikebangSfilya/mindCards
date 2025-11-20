package dtoout

import "time"

type DTOdel_out struct {
	Status  string    `json:"Status"`
	Object  string    `json:"Object"`
	TimeDel time.Time `json:"TimeDel"`
}

type MindCardDTO struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Tag         string    `json:"tag"`
	CreatedAt   time.Time `json:"created_at"`
	LevelStudy  int8      `json:"level_study"`
	Learned     bool      `json:"learned"`
}

func NewDelDTO(obj string) DTOdel_out {
	return DTOdel_out{
		Status:  "deleted",
		Object:  obj,
		TimeDel: time.Now(),
	}
}
