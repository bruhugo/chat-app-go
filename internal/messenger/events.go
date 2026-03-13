package messenger

import (
	"time"

	"github.com/google/uuid"
	"github.com/grongoglongo/chatter-go/internal/models"
	"github.com/grongoglongo/chatter-go/internal/models/dto"
)

type EventType string

const (
	CREATE_MESSAGE_EVENT_TYPE EventType = "CREATE_MESSAGE"
	DELETE_MESSAGE_EVENT_TYPE EventType = "DELETE_MESSAGE"
	LEAVE_CHAT_EVENT_TYPE     EventType = "LEAVE_CHAT"
	ENTER_CHAT_EVENT_TYPE     EventType = "ENTER_CHAT"
	START_TYPING_EVENT_TYPE   EventType = "START_TYPING"
	STOP_TYPING_EVENT_TYPE    EventType = "STOP_TYPING"
)

type EventWrapper struct {
	EventType EventType   `json:"eventType"`
	Chat      dto.ChatDto `json:"chat"`
	Event     any         `json:"event"`
	EventId   string      `json:"eventId"`
}

func CreateEventWrapper(eventType EventType, chat models.Chat, event any) *EventWrapper {
	return &EventWrapper{
		Chat:      *chat.ToDto(),
		EventType: eventType,
		Event:     event,
		EventId:   uuid.NewString(),
	}
}

type CreateMessageEvent struct {
	MessageId int64       `json:"messageId"`
	Content   string      `json:"content"`
	User      dto.UserDto `json:"user"`
	CreatedAt time.Time   `json:"createdAt"`
}

type DeleteMessageEvent struct {
	MessageId int64 `json:"messageId"`
}

type LeaveChatEvent struct {
	User  models.User `json:"user"`
	Actor models.User `json:"actor"`
}

type EnterChatEvent struct {
	User  models.User `json:"user"`
	Actor models.User `json:"actor"`
}

type TypingEvent struct {
	User   dto.UserDto `json:"user"`
	Typing bool        `json:"typing"`
}
