package cards

import "errors"

var (
	ErrNotExist  = errors.New("card not exist")
	errFailToAdd = errors.New("failed to add card")
)
