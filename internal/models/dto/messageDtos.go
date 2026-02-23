package dto

type MessageDto struct {
	ID      int64    `json:"id"`
	Content string   `json:"content"`
	User    *UserDto `json:"user"`
	Chat    *ChatDto `json:"chat"`
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

type WritingEventDto struct {
	Typing bool `json:"typing"`
}
