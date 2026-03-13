package dto

import (
	"time"
)

type Role string

const (
	USER  Role = "USER"
	ADMIN Role = "ADMIN"
)

type ChatDto struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	Creator     *UserDto
}

type CreateChatDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateChatDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ChatMemberDto struct {
	ID   int64   `json:"id"`
	Role Role    `json:"role"`
	Chat ChatDto `json:"chat"`
	User UserDto `json:"user"`
}

type AddChatMemberDto struct {
	Username string `json:"username"`
	Role     Role   `json:"role"`
}

type ChatResponseDto struct {
	ChatDto     *ChatDto    `json:"chat"`
	LastMessage *MessageDto `json:"lastMessage"`
}

type ChangeRoleDto struct {
	NewRole  Role  `json:"newRole"`
	TargetId int64 `json:"targetId"`
}

type DeleteChatMemberDto struct {
	TargetId int64 `json:"targetId"`
}
