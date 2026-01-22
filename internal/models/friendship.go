package models

type Friendship struct {
	ID        int64 `db:"id"`
	User1     User
	User2     User
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}
