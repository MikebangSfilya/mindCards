package cards

import (
	"fmt"
	"time"
)

var errAllFieldNeeder = fmt.Errorf("all fields are required")

type MindCard struct {
	ID          int64     `json:"id" db:"card_id"`
	UserID      int64     `json:"user_id" db:"card_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"card_description"`
	Tag         string    `json:"tag" db:"tag"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	LevelStudy  int8      `json:"level_study" db:"level_study"`
	Learned     bool      `json:"learned" db:"learned"`
}

func NewCard(title, description, tag string) (*MindCard, error) {
	return &MindCard{
		Title:       title,
		Description: description,
		Tag:         tag,
		CreatedAt:   time.Now(),
		LevelStudy:  0,
		Learned:     false,
	}, nil

}
