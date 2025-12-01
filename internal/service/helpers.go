package service

import (
	"cards/internal/model"
	"cards/internal/storage"
)

func rowsToCard(rows []storage.CardRow) []model.MindCard {

	result := make([]model.MindCard, len(rows))

	for i, row := range rows {
		result[i] = rowToCard(row)
	}

	return result

}

func rowToCard(row storage.CardRow) model.MindCard {
	return model.MindCard{
		ID:          row.ID,
		Title:       row.Title,
		Description: row.Description,
		Tag:         row.Tag,
		CreatedAt:   row.CreatedAt,
		LevelStudy:  row.LevelStudy,
		Learned:     row.Learned,
	}
}
