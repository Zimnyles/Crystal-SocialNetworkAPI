package friends

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type FriendsRepository struct {
	Dbpool       *pgxpool.Pool
	CustomLogger *zerolog.Logger
}

func NewFriendsRepository(dbpool *pgxpool.Pool, customLogger *zerolog.Logger) *FriendsRepository {
	return &FriendsRepository{
		Dbpool:       dbpool,
		CustomLogger: customLogger,
	}
}
