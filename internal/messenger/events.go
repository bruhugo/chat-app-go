package messenger

import (
	"time"
)

type EventType string

const (
	CREATE_MESSAGE_EVENT_TYPE EventType = "CREATE_MESSAGE"
	DELETE_MESSAGE_EVENT_TYPE EventType = "UPDATE_MESSAGE"
	LEAVE_CHAT_EVENT_TYPE     EventType = "LEAVE_CHAT"
	ENTER_CHAT_EVENT_TYPE     EventType = "ENTER_CHAT"
)

type EventWrapper struct {
	EventType EventType
	ChatId    int64
	Event     any
}

type CreateMessageEvent struct {
	MessageId int64
	Content   string
	UserId    int64
	CreatedAt time.Time
}

type DeleteMessageEvent struct {
	MessageId int64
}

type LeaveChatEvent struct {
	UserId  int64
	ActorId int64
}

type EnterChatEvent struct {
	UserId  int64
	ActorId int64
}
