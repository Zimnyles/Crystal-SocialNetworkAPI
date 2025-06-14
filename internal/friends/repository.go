package friends

import (
	"context"
	"fmt"
	"zimniyles/fibergo/internal/models"

	"github.com/jackc/pgx/v5"
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

func (r *FriendsRepository) GetAllFriendRequests(userID int) ([]models.FriendRequestList, error) {
    query := `
        SELECT 
            u.login,
            u.avatarpath,
            CASE
                WHEN f.user_id = @userID THEN 'pendingOutgoing'
                WHEN f.friend_id = @userID THEN 'pendingIncoming'
            END as status
        FROM friends f
        JOIN users u ON 
            (f.user_id = @userID AND f.friend_id = u.id) OR
            (f.friend_id = @userID AND f.user_id = u.id)
        WHERE f.status = 'pending'`

    args := pgx.NamedArgs{
        "userID": userID,
    }

    rows, err := r.Dbpool.Query(context.Background(), query, args)
    if err != nil {
        return nil, fmt.Errorf("error querying all friend requests: %w", err)
    }
    defer rows.Close()

    var requests []models.FriendRequestList
    for rows.Next() {
        var req models.FriendRequestList
        if err := rows.Scan(
            &req.Login,
            &req.AvatarPath,
            &req.FriendshipStatus,
        ); err != nil {
            return nil, fmt.Errorf("error scanning friend request: %w", err)
        }
        requests = append(requests, req)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error during rows iteration: %w", err)
    }

    return requests, nil
}

func (r *FriendsRepository) GetAcceptedFriends(userID int) ([]models.FriendList, error) {
    query := `
        SELECT 
            u.login,
            u.avatarpath,
            f.status
        FROM friends f
        JOIN users u ON 
            (f.user_id = @userID AND f.friend_id = u.id) OR
            (f.friend_id = @userID AND f.user_id = u.id)
        WHERE f.status = 'accepted'`

    args := pgx.NamedArgs{
        "userID": userID,
    }

    rows, err := r.Dbpool.Query(context.Background(), query, args)
    if err != nil {
        return nil, fmt.Errorf("error querying accepted friends: %w", err)
    }
    defer rows.Close()

    var friends []models.FriendList
    for rows.Next() {
        var friend models.FriendList
        if err := rows.Scan(
            &friend.Login,
            &friend.AvatarPath,
            &friend.FriendshipStatus,
        ); err != nil {
            return nil, fmt.Errorf("error scanning friend data: %w", err)
        }
        friends = append(friends, friend)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error during rows iteration: %w", err)
    }

    return friends, nil
}

func (r *FriendsRepository) GetIDfromLogin(login string) (int, error) {
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