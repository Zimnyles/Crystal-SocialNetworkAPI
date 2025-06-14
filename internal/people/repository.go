package people

import (
	"context"
	"fmt"
	"time"
	"zimniyles/fibergo/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type PeopleRepository struct {
	Dbpool       *pgxpool.Pool
	CustomLogger *zerolog.Logger
}

type friendshipStatus struct {
	Id        int
	Status    string
	CreatedAt time.Time
    OriginatorID int
    RecipientID int
}

func NewPeopleRepository(dbpool *pgxpool.Pool, customLogger *zerolog.Logger) *PeopleRepository {
	return &PeopleRepository{
		Dbpool:       dbpool,
		CustomLogger: customLogger,
	}
}

func (r *PeopleRepository) AddFriend(userID int, friendID int) (*friendshipStatus, error) {
	query := `
    INSERT INTO friends (user_id, friend_id, status)
    VALUES (@user_id, @friend_id, 'pending')
    ON CONFLICT (user_id, friend_id) 
    DO UPDATE SET 
        status = EXCLUDED.status,
        updated_at = CURRENT_TIMESTAMP
    RETURNING id, status, created_at`

	args := pgx.NamedArgs{
		"user_id":   userID,
		"friend_id": friendID,
	}

	var friendshipStatus friendshipStatus

	err := r.Dbpool.QueryRow(context.Background(), query, args).Scan(
		&friendshipStatus.Id,
		&friendshipStatus.Status,
		&friendshipStatus.CreatedAt,
	)

    friendshipStatus.OriginatorID = userID
    friendshipStatus.RecipientID = friendID


	if err != nil {
		return nil, fmt.Errorf("error creating friendship: %w", err)
	}

	return &friendshipStatus, nil

}

func (r *PeopleRepository) GetIDfromLogin(login string) (int, error) {
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

func (r *PeopleRepository) IsLoginExists(login string, logger *zerolog.Logger) (bool, error) {
	var exists bool
	err := r.Dbpool.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)",
		login,
	).Scan(&exists)

	return exists, err
}

func (r *PeopleRepository) CountAll() int {
	query := "SELECT count(*) FROM users"
	var count int
	r.Dbpool.QueryRow(context.Background(), query).Scan(&count)
	return count

}

func (r *PeopleRepository) CountNonFriends(currentUserID int) (int, error) {
    query := `
        SELECT COUNT(*) 
        FROM users u
        WHERE 
            u.id != @currentUserID  -- Исключаем текущего пользователя
            AND NOT EXISTS (
                
                SELECT 1 FROM friends 
                WHERE (user_id = @currentUserID AND friend_id = u.id AND status = 'accepted')
                   OR (user_id = u.id AND friend_id = @currentUserID AND status = 'accepted')
            )
            AND NOT EXISTS (
                
                SELECT 1 FROM friends 
                WHERE (user_id = @currentUserID AND friend_id = u.id AND status = 'blocked')
                   OR (user_id = u.id AND friend_id = @currentUserID AND status = 'blocked')
            )`

    var count int
    err := r.Dbpool.QueryRow(
        context.Background(),
        query,
        pgx.NamedArgs{"currentUserID": currentUserID},
    ).Scan(&count)

    if err != nil {
        return 0, fmt.Errorf("failed to count non-friends: %w", err)
    }

    return count, nil
}

func (r *PeopleRepository) GetAll(currentUserID int, limit, offset int, searchTerm string) ([]models.PeopleProfileCredentials, error) {
    query := `
        SELECT 
            u.login,
            u.avatarpath,
            u.role,
            false AS is_friend_to_user  
        FROM 
            users u
        WHERE 
            u.id != @currentUserID  
            AND NOT EXISTS (
                SELECT 1 FROM friends 
                WHERE (user_id = @currentUserID AND friend_id = u.id)
                   OR (user_id = u.id AND friend_id = @currentUserID)
            )`

    args := pgx.NamedArgs{
        "currentUserID": currentUserID,
        "limit":        limit,
        "offset":       offset,
    }

    if searchTerm != "" {
        query += ` AND u.login ILIKE '%' || @searchTerm || '%'`
        args["searchTerm"] = searchTerm
    }

    query += `
        ORDER BY 
            u.login ASC
        LIMIT @limit OFFSET @offset`

    rows, err := r.Dbpool.Query(context.Background(), query, args)
    if err != nil {
        return nil, fmt.Errorf("error querying non-friends: %w", err)
    }
    defer rows.Close()

    var users []models.PeopleProfileCredentials
    for rows.Next() {
        var user models.PeopleProfileCredentials
        err := rows.Scan(
            &user.Login,
            &user.AvatarPath,
            &user.Role,
            &user.IsFriendToUser,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning user row: %w", err)
        }
        users = append(users, user)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error during rows iteration: %w", err)
    }

    return users, nil
}