package dto

type UpdateMessageDto struct {
	NewContent string
	SenderId   int64
	MessageId  int64
}
