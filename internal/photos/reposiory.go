package photos

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type PhotosRepository struct {
	Dbpool       *pgxpool.Pool
	CustomLogger *zerolog.Logger
}

func NewPhotosRepository(dbpool *pgxpool.Pool, customLogger *zerolog.Logger) *PhotosRepository {
	return &PhotosRepository{
		Dbpool:       dbpool,
		CustomLogger: customLogger,
	}
}
