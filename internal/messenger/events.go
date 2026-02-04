package messenger

import "time"

type EventType string

const (
	CREATE_MESSAGE_EVENT_TYPE EventType = "CREATE_MESSAGE"
	DELETE_MESSAGE_EVENT_TYPE EventType = "UPDATE_MESSAGE"
)

type Event interface {
	EventType() EventType
}

type CreateMessageEvent struct {
	MessageId int64
	Content   string
	ChatId    int64
	UserId    int64
	CreatedAt time.Time
}

func (CreateMessageEvent) EventType() EventType {
	return CREATE_MESSAGE_EVENT_TYPE
}

type DeleteMessageEvent struct {
	MessageId int64
	ChatId    int64
}

func (DeleteMessageEvent) EventType() EventType {
	return DELETE_MESSAGE_EVENT_TYPE
}
