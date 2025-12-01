package cards

import "errors"

var (
	ErrNotExist   = errors.New("card not exist")
	errFailToAdd  = errors.New("failed to add card")
	ErrDecodeJSON = errors.New("failed to decode JSON")
	ErrEncodeJSON = errors.New("failed to encode response")
	ErrAddCard    = errors.New("failed to add card")
	ErrDeleteCard = errors.New("failed to delete card")
	ErrUpdateCard = errors.New("failed to update card")
)
