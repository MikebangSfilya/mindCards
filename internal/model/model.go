package model

import (
	"time"
)

type MindCard struct {
	title       string
	description string
	tag         string
	createdAt   time.Time
	levelStudy  int8
	learned     bool
}

func NewCard(title, description, tag string) *MindCard {
	return &MindCard{
		title:       title,
		description: description,
		tag:         tag,
		createdAt:   time.Now(),
		levelStudy:  0,
		learned:     false,
	}

}
