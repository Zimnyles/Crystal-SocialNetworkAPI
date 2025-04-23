package feed

import (
	"context"

	"github.com/jackc/pgx/v5"
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

func (r *FeedRepository) AddNewFeedPost(login string, content string, imagePath string) error {
	query := `INSERT INTO feedposts (login, image_path, content) 
          VALUES (@login, @image_path, @content)`
	// <input id="imageInput" type="file" name="image" accept="image/*" >
	args := pgx.NamedArgs{

		"login":      login,
		"image_path": imagePath,
		"content":    content,
	}

	_, err := r.Dbpool.Exec(context.Background(), query, args)
	if err != nil {
		r.CustomLogger.Info().Msg("не удалость создать пост")
	}
	return nil
}
