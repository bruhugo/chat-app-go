package messenger

import (
	"time"

	"github.com/grongoglongo/chatter-go/internal/models"
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
	Chat      models.Chat
	Event     any
}

type CreateMessageEvent struct {
	MessageId int64
	Content   string
	User      models.User
	CreatedAt time.Time
}

type DeleteMessageEvent struct {
	MessageId int64
}

type LeaveChatEvent struct {
	User  models.User
	Actor models.User
}

type EnterChatEvent struct {
	User  models.User
	Actor models.User
}
