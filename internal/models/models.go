package models

import "time"

type ProfileCredentials struct {
	Login      string    `db:"login"`
	Email      string    `db:"email"`
	Createdat  time.Time `db:"createdat"`
	Role       string    `db:"role"`
	AvatarPath string    `db:"avatarpath"`
}

type FeedPost struct {
	Id                string    `db:"id"`
	CreatorLogin      string    `db:"login"`
	Content           string    `db:"content"`
	ImagePath         string    `db:"image_path"`
	CreatedAt         time.Time `db:"created_at"`
	CreatorAvatarPath string    `db:"avatarpath"`
}

type PeopleProfileCredentials struct {
	Login          string `db:"login"`
	AvatarPath     string `db:"avatarpath"`
	Role           int    `db:"role"`
	IsFriendToUser bool
}

type FriendPageCredentials struct {
    Friends         []FriendList         
    FriendRequests  []FriendRequestList 
}

type FriendList struct {
	Login            string `db:"login"`
	AvatarPath       string `db:"avatarpath"`
	FriendshipStatus string `db:"status"`
	Role int `db:"role"`
}

type FriendRequestList struct {
	Login            string `db:"login"`
	AvatarPath       string `db:"avatarpath"`
	FriendshipStatus string `db:"status"`
}

type Message struct {
	Sender    string
	Receiver  string
	Content   string
	Timestamp time.Time
}

type Chat struct {
	Messages []Message
	Users    map[string]bool
}

func NewChat() *Chat {
    return &Chat{
        Messages: make([]Message, 0),
        Users:    make(map[string]bool),
    }
}

func (c *Chat) AddMessage(msg Message) {
    msg.Timestamp = time.Now()
    c.Messages = append(c.Messages, msg)
}

func (c *Chat) AddUser(username string) {
    c.Users[username] = true
}
