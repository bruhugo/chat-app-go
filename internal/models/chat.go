package models

import (
	"time"

	"github.com/grongoglongo/chatter-go/internal/models/dto"
)

type Chat struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	Creator     *User
}

func (c *Chat) ToDto() *dto.ChatDto {
	return &dto.ChatDto{
		Name:        c.Name,
		Description: c.Description,
		ID:          c.ID,
		CreatedAt:   c.CreatedAt,
		Creator: &dto.UserDto{
			Email:    c.Creator.Email,
			ID:       c.Creator.ID,
			Username: c.Creator.Username,
		},
	}
}
