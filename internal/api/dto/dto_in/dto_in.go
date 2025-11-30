package dtoin

// The Card is a descriptoin of mindCard DTO. Collects data to create card
type Card struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Tag         string `json:"tag"`
}

type Update struct {
	NewDeccription string `json:"description"`
}
