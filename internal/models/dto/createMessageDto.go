package dto

type CreateMessageDto struct {
	Content      string `json:"content"`
	TargetUserId int64  `json:"targetUserId"`
}
