package model

import (
	"testing"
)

func TestNewMode(t *testing.T) {
	testInput := struct {
		title       string
		description string
		tag         string
	}{
		title:       "TestName",
		description: "Test description in test case",
		tag:         "#Test",
	}

	testCard := NewCard(testInput.title, testInput.description, testInput.tag)
	if testCard == nil {
		t.Fatal("expected card struct, got nil")
	}

	switch {
	case testCard.Title != testInput.title:
		t.Errorf("extected %s, got %s", testInput.title, testCard.Title)
	case testCard.Description != testInput.description:
		t.Errorf("expected %s, gos %s", testCard.Description, testInput.description)
	case testCard.Tag != testInput.tag:
		t.Errorf("expected %s, gos %s", testCard.Tag, testInput.tag)
	default:
		t.Log("all tests passed")
	}

}
