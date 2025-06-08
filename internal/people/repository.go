package people

import (
	"context"
	"zimniyles/fibergo/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type PeopleRepository struct {
	Dbpool       *pgxpool.Pool
	CustomLogger *zerolog.Logger
}

func NewPeopleRepository(dbpool *pgxpool.Pool, customLogger *zerolog.Logger) *PeopleRepository {
	return &PeopleRepository{
		Dbpool:       dbpool,
		CustomLogger: customLogger,
	}
}

func (r *PeopleRepository) CountAll() int {
	query := "SELECT count(*) FROM users"
	var count int
	r.Dbpool.QueryRow(context.Background(), query).Scan(&count)
	return count

}

func (r *PeopleRepository) GetAll(limit, offset int) ([]models.PeopleProfileCredentials, error) {
	query := `
    SELECT 
        u.login,
        u.avatarpath,
        u.role
    FROM 
        users u
    ORDER BY 
        u.createdat DESC
	LIMIT @limit OFFSET @offset
	`
	args := pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	}
	rows, err := r.Dbpool.Query(context.Background(), query, args)

	if err != nil {
		return nil, err
	}
	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.PeopleProfileCredentials])
	if err != nil {
		return nil, err
	}
	return users, nil
}
