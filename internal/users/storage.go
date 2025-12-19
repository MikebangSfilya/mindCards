package users

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersStorage struct {
	db *pgxpool.Pool
}

func NewUserPool(db *pgxpool.Pool) *UsersStorage {
	repo := &UsersStorage{
		db: db,
	}

	return repo
}

func (s *UsersStorage) SaveUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users
		(email, password_hash)
		VALUES ($1, $2)
		RETURNING user_id
	`

	err := s.db.QueryRow(ctx, query, user.Email, user.EncryptedPassword).Scan(&user.UserId)
	if err != nil {
		return err
	}

	return nil
}
