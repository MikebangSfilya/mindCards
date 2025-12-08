package cards

import (
	"time"
)

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

func NewCard(title, description, tag string) *MindCard {
	return &MindCard{
		Title:       title,
		Description: description,
		Tag:         tag,
		CreatedAt:   time.Now(),
		LevelStudy:  0,
		Learned:     false,
	}

}
