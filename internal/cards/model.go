package cards

import (
	"fmt"
	"time"
)

var errAllFieldNeeder = fmt.Errorf("all fields are required")

type MindCard struct {
	CardID      int64     `json:"card_id"`
	UserID      int64     `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Tag         string    `json:"tag"`
	CreatedAt   time.Time `json:"created_at"`
	LevelStudy  int8      `json:"level_study"`
	Learned     bool      `json:"learned"`
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
