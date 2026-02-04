package dto

type GetMessagesDto struct {
	PageRequest PageRequest
	SenderId    int64
	ReceiverId  int64
}
