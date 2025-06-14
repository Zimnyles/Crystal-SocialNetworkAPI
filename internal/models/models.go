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
}

type FriendRequestList struct {
	Login            string `db:"login"`
	AvatarPath       string `db:"avatarpath"`
	FriendshipStatus string `db:"status"`
}
