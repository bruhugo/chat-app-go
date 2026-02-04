package dto

type MessageDto struct {
	ID       int64  `db:"id"`
	Content  string `db:"content"`
	Sender   *UserDto
	Receiver *UserDto
}

type GetMessagesDto struct {
	PageRequest PageRequest
	SenderId    int64
	ReceiverId  int64
}

type CreateMessageDto struct {
	Content      string `json:"content"`
	TargetUserId int64  `json:"targetUserId"`
}

type UpdateMessageDto struct {
	NewContent string
	SenderId   int64
	MessageId  int64
}
