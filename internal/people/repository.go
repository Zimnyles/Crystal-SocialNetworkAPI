package people

import (
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
