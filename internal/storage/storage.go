package storage

import "time"

type CardRow struct {
	ID          int64
	Title       string
	Description string
	Tag         string
	CreatedAt   time.Time
	LevelStudy  int8
	Learned     bool
}
