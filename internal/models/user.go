package models

type User struct {
	ID        int64  `db:"id"`
	Username  string `db:"username"`
	Email     string `db:"email"`
	Password  string `db:"password_hash"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}
