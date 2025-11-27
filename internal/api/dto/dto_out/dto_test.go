package dtoout

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDelTro(t *testing.T) {
	obj := "Test"
	delDto := NewDelDTO(obj)
	require.NotNil(t, delDto)
	require.Equal(t, "deleted", delDto.Status)
	require.Equal(t, "Test", delDto.Object)

}

func TestNewErr(t *testing.T) {
	var ErrTest = errors.New("err to test")
	ErrDto := NewErr(ErrTest)
	require.NotNil(t, ErrDto)
	require.Equal(t, "err to test", ErrDto.Err)
}

func TestNewErrToString(t *testing.T) {
	var ErrTest = errors.New("err to test")
	ErrDto := NewErr(ErrTest)
	require.NotNil(t, ErrDto)

	stringDto := ErrDto.ToString()
	require.NotNil(t, stringDto)

	var decoded map[string]interface{}
	err := json.Unmarshal([]byte(stringDto), &decoded)
	require.NoError(t, err)

	require.Contains(t, decoded, "Err")
	require.Equal(t, "err to test", decoded["Err"])

}
