package dto

import "time"

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
