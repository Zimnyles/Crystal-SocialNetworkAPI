package models

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type GlobalRepo interface {
	GetIDfromLogin(login string) (int, error)
}

type GlobalRepository struct {
	Dbpool       *pgxpool.Pool
	CustomLogger *zerolog.Logger
}

func NewGlobalRepository(dbpool *pgxpool.Pool, customLogger *zerolog.Logger) *GlobalRepository {
	return &GlobalRepository{
		Dbpool:       dbpool,
		CustomLogger: customLogger,
	}
}

func (r *GlobalRepository) GetIDfromLogin(login string) (int, error) {
	query := `
        SELECT id 
        FROM users 
        WHERE login = @login`

	args := pgx.NamedArgs{
		"login": login,
	}

	var userId int
	err := r.Dbpool.QueryRow(context.Background(), query, args).Scan(&userId)
	if err != nil {
		return 0, fmt.Errorf("cannot get user ID, server error: %w", err)
	}
	return userId, nil

}
