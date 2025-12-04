package cards

import (
	"testing"
	"time"

	"github.com/MikebangSfilya/mindCards/internal/storage"

	"github.com/stretchr/testify/require"
)

func TestRowToCard(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name     string
		row      storage.CardRow
		expected MindCard
	}{
		{
			name: "normal_card",
			row: storage.CardRow{
				CardID:      1,
				UserID:      10,
				Title:       "Go Interfaces",
				Description: "An interface type is defined as a set of method signatures",
				Tag:         "golang",
				CreatedAt:   now,
				LevelStudy:  2,
				Learned:     false,
			},
			expected: MindCard{
				CardID:      1,
				UserID:      10,
				Title:       "Go Interfaces",
				Description: "An interface type is defined as a set of method signatures",
				Tag:         "golang",
				CreatedAt:   now,
				LevelStudy:  2,
				Learned:     false,
			},
		},
		{
			name: "zero_values",
			row: storage.CardRow{
				CardID:      0,
				UserID:      10,
				Title:       "",
				Description: "",
				Tag:         "",
				CreatedAt:   time.Time{},
				LevelStudy:  0,
				Learned:     false,
			},
			expected: MindCard{
				CardID:      0,
				UserID:      10,
				Title:       "",
				Description: "",
				Tag:         "",
				CreatedAt:   time.Time{},
				LevelStudy:  0,
				Learned:     false,
			},
		},
		{
			name: "learned_card",
			row: storage.CardRow{
				CardID:      999,
				Title:       "Max Level",
				Description: "Already learned",
				Tag:         "done",
				CreatedAt:   now.Add(-24 * time.Hour),
				LevelStudy:  5,
				Learned:     true,
			},
			expected: MindCard{
				CardID:      999,
				Title:       "Max Level",
				Description: "Already learned",
				Tag:         "done",
				CreatedAt:   now.Add(-24 * time.Hour),
				LevelStudy:  5,
				Learned:     true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := rowToCard(tc.row)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestRowsToCards(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name     string
		rows     []storage.CardRow
		expected []MindCard
	}{
		{
			name: "multiple_cards",
			rows: []storage.CardRow{
				{
					CardID:      1,
					Title:       "First",
					Description: "Desc 1",
					Tag:         "tag1",
					CreatedAt:   now,
					LevelStudy:  0,
					Learned:     false,
				},
				{
					CardID:      2,
					Title:       "Second",
					Description: "Desc 2",
					Tag:         "tag2",
					CreatedAt:   now.Add(time.Hour),
					LevelStudy:  3,
					Learned:     true,
				},
			},
			expected: []MindCard{
				{
					CardID:      1,
					Title:       "First",
					Description: "Desc 1",
					Tag:         "tag1",
					CreatedAt:   now,
					LevelStudy:  0,
					Learned:     false,
				},
				{
					CardID:      2,
					Title:       "Second",
					Description: "Desc 2",
					Tag:         "tag2",
					CreatedAt:   now.Add(time.Hour),
					LevelStudy:  3,
					Learned:     true,
				},
			},
		},
		{
			name:     "empty_slice",
			rows:     []storage.CardRow{},
			expected: []MindCard{},
		},
		{
			name:     "nil_slice",
			rows:     nil,
			expected: []MindCard{},
		},
		{
			name: "single_card",
			rows: []storage.CardRow{
				{
					CardID:      42,
					Title:       "Single",
					Description: "Just one",
					Tag:         "alone",
					CreatedAt:   now,
					LevelStudy:  1,
					Learned:     false,
				},
			},
			expected: []MindCard{
				{
					CardID:      42,
					Title:       "Single",
					Description: "Just one",
					Tag:         "alone",
					CreatedAt:   now,
					LevelStudy:  1,
					Learned:     false,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := rowsToCards(tc.rows)
			require.Equal(t, tc.expected, result)
		})
	}
}
