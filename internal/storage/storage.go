package storage

import "time"

type CardRow struct {
	CardID      int64
	UserID      int64
	Title       string
	Description string
	Tag         string
	CreatedAt   time.Time
	LevelStudy  int8
	Learned     bool
}
