package models

import "github.com/grongoglongo/chatter-go/internal/models/dto"

type ChatMember struct {
	ID   int64    `db:"id"`
	Role dto.Role `db:"role"`
	Chat Chat
	User User
}

func (c ChatMember) ToDto() *dto.ChatMemberDto {
	return &dto.ChatMemberDto{
		ID:   c.ID,
		Chat: *c.Chat.ToDto(),
		User: *c.User.ToDto(),
		Role: c.Role,
	}
}
