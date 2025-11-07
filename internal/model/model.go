package model

import (
	"time"
)

type MindCard struct {
	Title       string
	Description string
	Tag         string
	CreatedAt   time.Time
	LevelStudy  int8
	Learned     bool
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
