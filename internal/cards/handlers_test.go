package cards

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)
	require.NotNil(t, dtoIn)
	assert.Equal(t, dto.Title, dtoIn.Title)
	assert.Equal(t, dto.Description, dtoIn.Description)

	assert.Equal(t, dto.Tag, dtoIn.Tag)
}

func TestBase(t *testing.T) {
}
