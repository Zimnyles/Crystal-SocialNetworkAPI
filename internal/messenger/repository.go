package messenger

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type MessengerRepository struct {
	Dbpool       *pgxpool.Pool
	CustomLogger *zerolog.Logger
}

func NewMessengerRepository(dbpool *pgxpool.Pool, customLogger *zerolog.Logger) *MessengerRepository {
	return &MessengerRepository{
		Dbpool:       dbpool,
		CustomLogger: customLogger,
	}
}
