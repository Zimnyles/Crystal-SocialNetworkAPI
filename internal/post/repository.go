package post

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type PostRepository struct {
	Dbpool       *pgxpool.Pool
	CustomLogger *zerolog.Logger
}

func NewPostRepository(dbpool *pgxpool.Pool, customLogger *zerolog.Logger) *PostRepository {
	return &PostRepository{
		Dbpool:       dbpool,
		CustomLogger: customLogger,
	}
}

func (r *PostRepository) CountAll() int{
	query := "SELECT count(*) FROM posts"
	var count int
	r.Dbpool.QueryRow(context.Background(), query).Scan(&count)
	return count

}

func (r *PostRepository) GetAll(limit, offset int) ([]Post, error) {
	query := "SELECT * from posts ORDER BY createdat LIMIT @limit OFFSET @offset"
	args := pgx.NamedArgs{
		"limit": limit,
		"offset": offset,
	}
	rows, err := r.Dbpool.Query(context.Background(), query, args)
	if err != nil {
		return nil, err
	}
	posts, err := pgx.CollectRows(rows, pgx.RowToStructByName[Post])
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PostRepository) addPost(form PostCreateForm, logger *zerolog.Logger) (int, error) {
	query := `
		INSERT INTO posts (email, name, price, location, description, breed, createdat) 
		VALUES (@email, @name, @price, @location, @description, @breed, @createdat)
		RETURNING id
	`
	args := pgx.NamedArgs{
		"email":       form.Email,
		"name":        form.Name,
		"price":       form.Price,
		"location":    form.Location,
		"description": form.Description,
		"breed":       form.Breed,
		"createdat":   time.Now(),
	}

	var id int

	err := r.Dbpool.QueryRow(context.Background(), query, args).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("невозможно создать пост: %w", err)
	}
	logger.Info().Msg("Создан пост с id: " + strconv.Itoa(id))
	return id, nil
}
