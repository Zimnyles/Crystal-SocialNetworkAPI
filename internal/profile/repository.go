package profile

import (
	"context"
	"fmt"
	"strings"
	"zimniyles/fibergo/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type ProfileRepository struct {
	Dbpool       *pgxpool.Pool
	CustomLogger *zerolog.Logger
}

func NewProfileRepository(dbpool *pgxpool.Pool, customLogger *zerolog.Logger) *ProfileRepository {
	return &ProfileRepository{
		Dbpool:       dbpool,
		CustomLogger: customLogger,
	}
}

func (r *ProfileRepository) IsLoginExistsForString(login string, logger *zerolog.Logger) (bool, error) {
	var exists bool
	err := r.Dbpool.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)",
		login,
	).Scan(&exists)

	return exists, err
}

func (r *ProfileRepository) GetUserPhotos(userID int, logger *zerolog.Logger) ([]models.PhotoList, error) {
	query := `
        SELECT id, user_id, created_at, file_path, is_public
        FROM photos
        WHERE user_id = $1 AND is_public = TRUE
        ORDER BY created_at DESC
    `
	rows, err := r.Dbpool.Query(context.Background(), query, userID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to fetch public photos")
		return nil, err
	}

	defer rows.Close()

	var photoList []models.PhotoList

	for rows.Next() {
		var photo models.PhotoList
		err := rows.Scan(
			&photo.PhotoID,
			&photo.UserID,
			&photo.CreatedAt,
			&photo.ImagePath,
			&photo.IsPublic,
		)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to scan photo row")
			return nil, err
		}

		photoList = append(photoList, photo)
	}

	if err := rows.Err(); err != nil {
		logger.Error().Err(err).Msg("Error during rows iteration")
		return nil, err
	}

	return photoList, nil

}

func (r *ProfileRepository) GetUserDataFromLogin(login string, logger *zerolog.Logger) (*models.ProfileCredentials, error) {
	query := `
        SELECT 
            login,  
            email,
			createdat,
			role,
			avatarpath 
        FROM users 
        WHERE login = @login`
	args := pgx.NamedArgs{
		"login": login,
	}

	row := r.Dbpool.QueryRow(context.Background(), query, args)
	var ProfileCredentials models.ProfileCredentials

	err := row.Scan(&ProfileCredentials.Login, &ProfileCredentials.Email, &ProfileCredentials.Createdat, &ProfileCredentials.Role, &ProfileCredentials.AvatarPath)

	if ProfileCredentials.AvatarPath != "" && !strings.HasPrefix(ProfileCredentials.AvatarPath, "/") {
		ProfileCredentials.AvatarPath = "/" + ProfileCredentials.AvatarPath
	}

	if err != nil {
		logger.Info().Msg("2")
		return nil, fmt.Errorf("error scanning password s36 : %w", err)
	}
	return &ProfileCredentials, nil
}

func (r *ProfileRepository) UpdateUserAvatar(id int, path string) error {
	query := `UPDATE users SET avatarpath = @avatarpath WHERE id = @id`

	args := pgx.NamedArgs{
		"avatarpath": path,
		"id":      id,
	}

	_, err := r.Dbpool.Exec(context.Background(), query, args)
	if err != nil {
		return fmt.Errorf("невозможно обновить аватар пост: %w", err)
	}
	return nil
}

func (r *ProfileRepository) AddUserAvatarToPhotos(userID int, filePath string, isPublic bool) error {
	query := `
        INSERT INTO photos (user_id, file_path, is_public)
        VALUES (@user_id, @file_path, @is_public)
    `

    args := pgx.NamedArgs{
        "user_id":   userID,
        "file_path": filePath,
        "is_public": isPublic,
    }

    _, err := r.Dbpool.Exec(context.Background(), query, args)
    if err != nil {
        return fmt.Errorf("не удалось добавить фото: %w", err)
    }
    return nil
}

func (r *ProfileRepository) GetAllUserPosts(userLogin string, limit int, offset int) ([]models.FeedPost, error) {
	query := `
        SELECT 
            fp.*,
            u.avatarpath
        FROM 
            feedposts fp
        LEFT JOIN 
            users u ON fp.login = u.login
        WHERE 
            fp.login = @userLogin  -- Предполагая, что есть поле user_id в feedposts
        ORDER BY 
            fp.created_at DESC
        LIMIT @limit OFFSET @offset
    `
	args := pgx.NamedArgs{
		"userLogin": userLogin,
		"limit":     limit,
		"offset":    offset,
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

func (r *ProfileRepository) CountUserPosts(userLogin string) (int, error) {
	query := `
        SELECT COUNT(*) 
        FROM feedposts 
        WHERE login = $1
    `

	var count int
	err := r.Dbpool.QueryRow(context.Background(), query, userLogin).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count user posts: %w", err)
	}

	return count, nil
}

func (r *ProfileRepository) GetIDfromLogin(login string) (int, error) {
	query := `
        SELECT id 
        FROM users 
        WHERE login = @login`

	args := pgx.NamedArgs{
		"login": login,
	}

	var userId int
	err := r.Dbpool.QueryRow(context.Background(), query, args).Scan(&userId)
	if err != nil {
		return 0, fmt.Errorf("cannot get user ID, server error: %w", err)
	}
	return userId, nil

}