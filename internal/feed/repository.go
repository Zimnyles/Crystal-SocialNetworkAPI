package feed

import (
	"context"
	"zimniyles/fibergo/internal/models"

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

func (r *FeedRepository) CountAll() int{
	query := "SELECT count(*) FROM feedposts"
	var count int
	r.Dbpool.QueryRow(context.Background(), query).Scan(&count)
	return count

}


func (r *FeedRepository) GetAll(limit, offset int) ([]models.FeedPost, error) {
	query := `
        SELECT 
            fp.*,
            u.avatarpath
        FROM 
            feedposts fp
        LEFT JOIN 
            users u ON fp.login = u.login
        ORDER BY 
            fp.created_at DESC
        LIMIT @limit OFFSET @offset
    `
	args := pgx.NamedArgs{
		"limit": limit,
		"offset": offset,
	}
	rows, err := r.Dbpool.Query(context.Background(), query, args)
	
	if err != nil {
		return nil, err
	}
	posts, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.FeedPost])
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *FeedRepository) AddNewFeedPost(login string, content string, imagePath string) error {
	query := `INSERT INTO feedposts (login, image_path, content) 
          VALUES (@login, @image_path, @content)`
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

func (r *FeedRepository) AddNewFeedPostWithoutImage(login string, content string) error {
	query := `INSERT INTO feedposts (login, image_path, content) 
          VALUES (@login, @image_path, @content)`
	args := pgx.NamedArgs{

		"login":      login,
		"image_path": "",
		"content":    content,
	}

	_, err := r.Dbpool.Exec(context.Background(), query, args)
	if err != nil {
		r.CustomLogger.Info().Msg("не удалость создать пост")
	}
	return nil
}
