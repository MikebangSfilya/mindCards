package dtoin

import (
	"fmt"
	"strings"
)

var (
	errEmptyDescription = fmt.Errorf("new description cant be empty")
	errShortDescription = fmt.Errorf("description is too short")
)

const minDescriptionLength = 10

// The Card is a descriptoin of mindCard DTO. Collects data to create card
type Card struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Tag         string `json:"tag"`
}

type Update struct {
	NewDescription string `json:"description"`
}

func (u *Update) Validate() error {

	trimmedDesc := strings.TrimSpace(u.NewDescription)

	if trimmedDesc == "" {
		return errEmptyDescription
	}

	//maybe should delete this validation
	if len(trimmedDesc) < minDescriptionLength {
		return errShortDescription
	}
	return nil
}
