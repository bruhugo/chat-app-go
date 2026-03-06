package models

import (
	"time"

	"github.com/grongoglongo/chatter-go/internal/models/dto"
)

type Message struct {
	ID        int64  `db:"id"`
	Content   string `db:"content"`
	User      *User
	Chat      *Chat
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (m *Message) ToDto() *dto.MessageDto {
	return &dto.MessageDto{
		ID:        m.ID,
		Content:   m.Content,
		User:      m.User.ToDto(),
		Chat:      m.Chat.ToDto(),
		UpdatedAt: m.UpdatedAt,
		CreatedAt: m.CreatedAt,
	}
}
