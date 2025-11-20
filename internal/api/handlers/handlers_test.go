package handlers

import (
	"bytes"
	dtoin "cards/internal/api/dto/dto_in"
	"encoding/json"
	"net/http"
	"testing"
)

func TestDecoder(t *testing.T) {
	dto := dtoin.Card{
		Title:       "title",
		Description: "Desc",
		Tag:         "Go",
	}
	testCase, _ := json.Marshal(dto)
	req, _ := http.NewRequest("POST", "/", bytes.NewReader(testCase))
	dtoIn := dtoin.Card{}
	err := decoder(req, &dtoIn)
	if err != nil {
		t.Errorf("decoder failed %v", err)
	}
	if dtoIn.Title == "" {
		t.Errorf("empty title")
	}
	if dtoIn.Description == "" {
		t.Errorf("empty description")
	}
	if dtoIn.Tag == "" {
		t.Errorf("empty tag")
	}

}
