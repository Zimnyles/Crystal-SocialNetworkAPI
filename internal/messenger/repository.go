package messenger

import (
	"context"
	"fmt"
	"zimniyles/fibergo/internal/models"

	"github.com/jackc/pgx/v5"
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

func (r *MessengerRepository) GetUserChats(userID int) ([]models.ChatPreview, error) {
    query := `
        WITH last_messages AS (
            SELECT 
                chat_id,
                content,
                sender_id,
                created_at,
                ROW_NUMBER() OVER (PARTITION BY chat_id ORDER BY created_at DESC) as rn
            FROM messages
        )
        SELECT 
            c.id AS chat_id,
            CASE 
                WHEN c.user1 = @user_id THEN u2.id 
                ELSE u1.id 
            END AS interlocutor_id,
            CASE 
                WHEN c.user1 = @user_id THEN u2.login 
                ELSE u1.login 
            END AS interlocutor_login,
            CASE
                WHEN lm.sender_id = @user_id THEN 'Вы: ' || COALESCE(lm.content, '')
                ELSE COALESCE(lm.content, 'Нет сообщений')
            END AS last_message,
            COALESCE(lm.created_at, c.created_at) AS last_message_time,
            lm.sender_id = @user_id AS is_your_message,
            CASE
                WHEN c.user1 = @user_id THEN u2.avatarpath
                ELSE u1.avatarpath
            END AS interlocutor_avatar_path
        FROM chats c
        JOIN users u1 ON u1.id = c.user1
        JOIN users u2 ON u2.id = c.user2
        LEFT JOIN last_messages lm ON lm.chat_id = c.id AND lm.rn = 1
        WHERE c.user1 = @user_id OR c.user2 = @user_id
        ORDER BY last_message_time DESC
    `

    args := pgx.NamedArgs{
        "user_id": userID,
    }

    rows, err := r.Dbpool.Query(context.Background(), query, args)
    if err != nil {
        return nil, fmt.Errorf("failed to get user chats: %w", err)
    }
    defer rows.Close()

    var chats []models.ChatPreview
    for rows.Next() {
        var cp models.ChatPreview
        err := rows.Scan(
            &cp.ChatID,
            &cp.InterlocutorID,
            &cp.InterlocutorLogin,
            &cp.LastMessage,
            &cp.LastMessageTime,
            &cp.IsYourMessage,
            &cp.InterlocutorAvatarPath,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan chat preview: %w", err)
        }
        chats = append(chats, cp)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating chat preview rows: %w", err)
    }

    return chats, nil
}