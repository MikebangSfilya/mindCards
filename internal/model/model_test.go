package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewModel(t *testing.T) {
	testInput := []struct {
		name        string
		title       string
		description string
		tag         string
		wantErr     bool
	}{
		{
			name:        "create card",
			title:       "TestName",
			description: "Test description in test case",
			tag:         "#Test",
			wantErr:     false,
		},
		{
			name:        "empty fields",
			title:       "",
			description: "",
			tag:         "",
			wantErr:     true,
		},
		{
			name:        "empty title",
			title:       "",
			description: "Test description in test case",
			tag:         "#Test",
			wantErr:     true,
		},
		{
			name:        "empty description",
			title:       "title",
			description: "",
			tag:         "#Test",
			wantErr:     true,
		},
		{
			name:        "empty tag",
			title:       "title",
			description: "Test description in test case",
			tag:         "",
			wantErr:     true,
		},
	}

	for _, tCase := range testInput {
		t.Run(tCase.name, func(t *testing.T) {
			card, err := NewCard(tCase.title, tCase.description, tCase.tag)
			if tCase.wantErr {
				require.Error(t, err)
				require.EqualError(t, errAllFieldNeeder, err.Error())
			}
			if !tCase.wantErr {
				require.NoError(t, err)
				require.Equal(t, tCase.title, card.Title)
				require.Equal(t, tCase.description, card.Description)
				require.Equal(t, tCase.tag, card.Tag)
			}
		})
	}

}
