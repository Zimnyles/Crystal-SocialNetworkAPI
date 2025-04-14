package profile

import (
	"context"
	"fmt"
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
			role 
        FROM users 
        WHERE login = @login`
	args := pgx.NamedArgs{
		"login": login,
	}
	
	row := r.Dbpool.QueryRow(context.Background(), query, args)
	var ProfileCredentials models.ProfileCredentials
	
	err := row.Scan(&ProfileCredentials.Login, &ProfileCredentials.Email, &ProfileCredentials.Createdat, &ProfileCredentials.Role)
	if err != nil {
		
		return nil, fmt.Errorf("error scanning password s36 : %w", err)
	}
	return &ProfileCredentials, nil
}
