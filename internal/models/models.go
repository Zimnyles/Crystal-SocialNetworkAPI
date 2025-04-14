package models

import "time"

type ProfileCredentials struct {
	Login     string `db:"login"`
	Email     string `db:"email"`
	Createdat time.Time `db:"createdat"`
	Role      string `db:"role"`
}
