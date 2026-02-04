package dto

type MessageDto struct {
	ID       int64  `db:"id"`
	Content  string `db:"content"`
	Sender   *UserDto
	Receiver *UserDto
}
