package models

import (
	"time"

	"github.com/grongoglongo/chatter-go/internal/models/dto"
)

type Message struct {
	ID        int64  `db:"id"`
	Content   string `db:"content"`
	Sender    *User
	Receiver  *User
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (m *Message) ToDto() *dto.MessageDto {
	return &dto.MessageDto{
		ID:       m.ID,
		Content:  m.Content,
		Sender:   m.Sender.ToDto(),
		Receiver: m.Sender.ToDto(),
	}
}
