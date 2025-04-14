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

func (r *UsersRepository) addUser(form UserCreateForm, logger *zerolog.Logger) (string, error) {

	emailIsExists, err := r.IsEmailExists(form, logger)
	if emailIsExists {
		return "Аккаунт с таким email уже существует", fmt.Errorf("невозможно зарегестрировать аккаунт : %w", err)
	}

	loginIsExists, err := r.IsLoginExists(form, logger)
	if loginIsExists {
		return "Аккаунт с таким логином уже существует", fmt.Errorf("невозможно зарегестрировать аккаунт : %w", err)
	}

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

	_, err = r.Dbpool.Exec(context.Background(), query, args)
	if err != nil {
		return "Ошибка сервера, попробуйте позже", fmt.Errorf("невозможно зарегестрировать аккаунт : %w", err)
	}
	logger.Info().Msg("зарегестрирован аккаунт")
	return "Аккаунт зарегестрирован", nil

}

func (r *UsersRepository) IsEmailExists(form UserCreateForm, logger *zerolog.Logger) (bool, error) {
	var exists bool
	err := r.Dbpool.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)",
		form.Email,
	).Scan(&exists)

	return exists, err
}

func (r *UsersRepository) IsLoginExists(form UserCreateForm, logger *zerolog.Logger) (bool, error) {
	var exists bool
	err := r.Dbpool.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)",
		form.Login,
	).Scan(&exists)

	return exists, err
}
