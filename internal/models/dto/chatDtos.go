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
	CreatedAt   time.Time `json:"created_at"`
	Creator     *UserDto
}

type CreateChatDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatorId   int64  `json:"creator_id"`
}

type UpdateChatDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ChatMemberDto struct {
	ID   int64 `db:"id"`
	Role Role  `db:"role"`
	Chat ChatDto
	User UserDto
}

type AddChatMemberDto struct {
	TargetId int64 `json:"target_id"`
	Role     Role  `json:"role"`
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
