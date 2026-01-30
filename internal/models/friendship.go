package models

import "time"

type Friendship struct {
	ID        int64 `db:"id"`
	User1     User
	User2     User
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
