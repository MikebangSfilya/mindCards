package model

import (
	"fmt"
	"time"
)

type MindCard struct {
	ID          int64     `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Tag         string    `json:"tag" db:"tag"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	LevelStudy  int8      `json:"level_study" db:"level_study"`
	Learned     bool      `json:"learned" db:"learned"`
}

func NewCard(title, description, tag string) (*MindCard, error) {
	if title == "" || description == "" || tag == "" {
		return nil, fmt.Errorf("all fields are required")
	}
	return &MindCard{
		Title:       title,
		Description: description,
		Tag:         tag,
		CreatedAt:   time.Now(),
		LevelStudy:  0,
		Learned:     false,
	}, nil

}
