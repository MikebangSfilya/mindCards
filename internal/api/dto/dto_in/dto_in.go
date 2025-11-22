package dtoin

type DTO interface{}

// The Card is a descriptoin of mindCard DTO. Collects data to create card
type Card struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Tag         string `json:"tag"`
}

// The DTODel is a DTO model that has a title to delete card by its name
type DTODel struct {
	Title string `json:"title"`
}

type Update struct {
	NewDeccription string `json:"description"`
}

// type UpdateLvlLearn struct {
// 	Title string `json:"title"`
// 	Level int    `json:"level"`
// }

type LimitOffset struct {
	Limit  int16 `json:"limit"`
	Offset int16 `json:"offset"`
}

func (g *LimitOffset) PaginationDefault() {
	if g.Limit == 0 {
		g.Limit = 50
	}
	if g.Limit > 1000 {
		g.Limit = 1000
	}
	if g.Offset < 0 {
		g.Offset = 0
	}
}
