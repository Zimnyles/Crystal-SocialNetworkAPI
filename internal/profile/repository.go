package profile

import (
	"context"
	"fmt"
	"strings"
	"zimniyles/fibergo/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type ProfileRepository struct {
	Dbpool       *pgxpool.Pool
	CustomLogger *zerolog.Logger
}

func NewProfileRepository(dbpool *pgxpool.Pool, customLogger *zerolog.Logger) *ProfileRepository {
	return &ProfileRepository{
		Dbpool:       dbpool,
		CustomLogger: customLogger,
	}
}

func (r *ProfileRepository) IsLoginExistsForString(login string, logger *zerolog.Logger) (bool, error) {
	var exists bool
	err := r.Dbpool.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)",
		login,
	).Scan(&exists)

	return exists, err
}

func (r *ProfileRepository) GetUserDataFromLogin(login string, logger *zerolog.Logger) (*models.ProfileCredentials, error) {
	logger.Info().Msg("1")
	query := `
        SELECT 
            login,  
            email,
			createdat,
			role,
			avatarpath 
        FROM users 
        WHERE login = @login`
	args := pgx.NamedArgs{
		"login": login,
	}

	row := r.Dbpool.QueryRow(context.Background(), query, args)
	var ProfileCredentials models.ProfileCredentials

	err := row.Scan(&ProfileCredentials.Login, &ProfileCredentials.Email, &ProfileCredentials.Createdat, &ProfileCredentials.Role, &ProfileCredentials.AvatarPath)

	if ProfileCredentials.AvatarPath != "" && !strings.HasPrefix(ProfileCredentials.AvatarPath, "/") {
		ProfileCredentials.AvatarPath = "/" + ProfileCredentials.AvatarPath
	}

	if err != nil {
		logger.Info().Msg("2")
		return nil, fmt.Errorf("error scanning password s36 : %w", err)
	}
	return &ProfileCredentials, nil
}

func (r *ProfileRepository) UpdateUserAvatar(login string, path string) error {
	query := `UPDATE users SET avatarpath = @avatarpath WHERE login = @login`

	args := pgx.NamedArgs{
		"avatarpath": path,
		"login":      login,
	}

	_, err := r.Dbpool.Exec(context.Background(), query, args)
	if err != nil {
		return fmt.Errorf("невозможно обновить аватар пост: %w", err)
	}
	return nil
}
