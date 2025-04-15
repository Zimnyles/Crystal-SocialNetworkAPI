package announcements

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type AnnouncementsRepository struct {
	Dbpool       *pgxpool.Pool
	CustomLogger *zerolog.Logger
}

func NewAnnouncementsRepository(dbpool *pgxpool.Pool, customLogger *zerolog.Logger) *AnnouncementsRepository {
	return &AnnouncementsRepository{
		Dbpool:       dbpool,
		CustomLogger: customLogger,
	}
}
