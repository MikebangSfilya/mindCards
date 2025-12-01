package cards

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestDecoder(t *testing.T) {
	dto := Card{
		Title:       "title",
		Description: "Desc",
		Tag:         "Go",
	}
	testCase, _ := json.Marshal(dto)
	req, _ := http.NewRequest("POST", "/", bytes.NewReader(testCase))
	dtoIn := Card{}
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

func TestBase(t *testing.T) {

}
