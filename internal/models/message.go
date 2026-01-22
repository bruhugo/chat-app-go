package models

type Message struct {
	ID       int64  `db:"id"`
	Content  string `db:"content"`
	Sender   User
	Receiver User
}
