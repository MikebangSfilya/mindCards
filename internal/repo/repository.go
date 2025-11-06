package repo

import "context"

type Repo interface {
	AddCard(ctx context.Context)
	UptadeCardDescription(ctx context.Context)
	DeleteCard(ctx context.Context)
}
