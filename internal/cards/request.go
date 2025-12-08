package cards

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
	UserID      int
	Title       string `json:"title"`
	Description string `json:"description"`
	Tag         string `json:"tag"`
}

func (c *Card) Validate() error {
	trimmedTitle := strings.TrimSpace(c.Title)
	trimmedDesc := strings.TrimSpace(c.Description)
	trimmedTag := strings.TrimSpace(c.Tag)

	if trimmedTitle == "" {
		return fmt.Errorf("error validate, title is empty.")
	} else if trimmedDesc == "" {
		return fmt.Errorf("error validate, description is empty.")
	} else if trimmedTag == "" {
		return fmt.Errorf("error validate, tag is empty.")
	}
	return nil
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
