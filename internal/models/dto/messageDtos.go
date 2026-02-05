package dto

type MessageDto struct {
	ID      int64  `db:"id"`
	Content string `db:"content"`
	User    *UserDto
	Chat    *ChatDto
}

type CreateMessageDto struct {
	Content string `json:"content"`
	ChatId  int64  `json:"chat_id"`
}

type UpdateMessageDto struct {
	NewContent string
	UserId     int64
	MessageId  int64
}
