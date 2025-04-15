package feed

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type FeedRepository struct {
	Dbpool       *pgxpool.Pool
	CustomLogger *zerolog.Logger
}

func NewFeedRepository(dbpool *pgxpool.Pool, customLogger *zerolog.Logger) *FeedRepository {
	return &FeedRepository{
		Dbpool:       dbpool,
		CustomLogger: customLogger,
	}
}