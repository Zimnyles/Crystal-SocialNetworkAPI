package home

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type UsersRepository struct {
	Dbpool       *pgxpool.Pool
	CustomLogger *zerolog.Logger
}

func NewUsersRepository(dbpool *pgxpool.Pool, customLogger *zerolog.Logger) *UsersRepository {
	return &UsersRepository{
		Dbpool:       dbpool,
		CustomLogger: customLogger,
	}
}

func (r *UsersRepository) addUser(form UserCreateForm, logger *zerolog.Logger) error {
	query := `
		INSERT INTO users (email, login, password, createdat) 
		VALUES (@email, @login, @password, @createdat)
	`
	args := pgx.NamedArgs{
		"email":     form.Email,
		"login":     form.Login,
		"password":  form.Password,
		"createdat": time.Now(),
	}

	
	_, err := r.Dbpool.Exec(context.Background(), query, args)
	if err != nil {
		return fmt.Errorf("невозможно зарегестрировать аккаунт : %w", err)
	}
	logger.Info().Msg("зарегестрирован аккаунт")
	return nil

}
