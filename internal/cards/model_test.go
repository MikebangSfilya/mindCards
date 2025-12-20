package cards

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewModel(t *testing.T) {
	testInput := []struct {
		name        string
		title       string
		description string
		tag         string
	}{
		{
			name:        "create card",
			title:       "TestName",
			description: "Test description in test case",
			tag:         "#Test",
		},
		{
			name:        "space valid description",
			title:       "title",
			description: "           Test description in test case       ",
			tag:         "#Test",
		},
		{
			name:        "space valid tag",
			title:       "title",
			description: "Test description in test case",
			tag:         "      tagg      ",
		},
	}

	for _, tCase := range testInput {
		t.Run(tCase.name, func(t *testing.T) {
			card := NewCard(tCase.title, tCase.description, tCase.tag)
			assert.Equal(t, tCase.title, card.Title)
			assert.Equal(t, tCase.description, card.Description)
			assert.Equal(t, tCase.tag, card.Tag)
		})
	}
}
