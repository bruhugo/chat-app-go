package dto

import "time"

type MessageDto struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	User      *UserDto  `json:"user"`
	Chat      *ChatDto  `json:"chat"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateMessageDto struct {
	Content string `json:"content"`
}

type UpdateMessageDto struct {
	NewContent string
	UserId     int64
	MessageId  int64
}

type WritingEventDto struct {
	Typing bool `json:"typing"`
}
