package models

type Role string

const (
	USER  Role = "USER"
	ADMIN Role = "ADMIN"
)

type ChatMember struct {
	ID   int64 `db:"id"`
	Role Role  `db:"role"`
	Chat Chat
	User User
}
