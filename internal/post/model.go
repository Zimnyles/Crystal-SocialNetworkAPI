package post

import "time"

type PostCreateForm struct {
	Email       string
	Name        string
	Price       string
	Breed       string
	Description string
	Location    string
}

type Post struct {
	Id          string       `db:"id"`
	Email       string    `db:"email"`
	Name        string    `db:"name"`
	Price       string    `db:"price"`
	Breed       string    `db:"breed"`
	Description string    `db:"description"`
	Location    string    `db:"location"`
	CreatedAt   time.Time `db:"createdat"`
}
