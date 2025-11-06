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
	case testCard.title != testInput.title:
		t.Errorf("extected %s, got %s", testInput.title, testCard.title)
	case testCard.description != testInput.description:
		t.Errorf("expected %s, gos %s", testCard.description, testInput.description)
	case testCard.tag != testInput.tag:
		t.Errorf("expected %s, gos %s", testCard.tag, testInput.tag)
	default:
		t.Log("all tests passed")
	}

}
