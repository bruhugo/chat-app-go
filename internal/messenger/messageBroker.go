package messenger

import (
	"context"
	"log"

	"github.com/grongoglongo/chatter-go/internal/models"
)

type Messenger interface {
	Post(EventWrapper) error
	Listen(context.Context, ...func(e EventWrapper) error) error
}

type EventBus struct {
	hub       *ConnectionHub
	messenger Messenger
	channel   chan EventWrapper
}

func NewEventBus(m Messenger, h *ConnectionHub) *EventBus {
	ma := &EventBus{
		messenger: m,
		hub:       h,
	}

	ctx := context.Background()

	go func() {
		m.Listen(ctx, func(e EventWrapper) error {
			h.Broadcast(e)
			return nil
		})
	}()

	return ma
}

func (bus *EventBus) PostCreateMessageEvent(m models.Message) {
	ew := EventWrapper{
		ChatId:    m.Chat.ID,
		EventType: CREATE_MESSAGE_EVENT_TYPE,
		Event: CreateMessageEvent{
			Content:   m.Content,
			UserId:    m.User.ID,
			CreatedAt: m.CreatedAt,
		},
	}
	bus.messenger.Post(ew)
	log.Printf("Create message event of message %d sent to chat %d", m.ID, m.Chat.ID)
}

func (bus *EventBus) PostDeleteMessageEvent(messageId, chatId int64) {
	ew := EventWrapper{
		ChatId:    chatId,
		EventType: CREATE_MESSAGE_EVENT_TYPE,
		Event: DeleteMessageEvent{
			MessageId: messageId,
		},
	}

	bus.messenger.Post(ew)
	log.Printf("Delete message event of message %d sent to chat %d", messageId, chatId)
}

func (bus *EventBus) PostLeaveChatEvent(chatId, userId, actorId int64) {
	ew := EventWrapper{
		ChatId:    chatId,
		EventType: LEAVE_CHAT_EVENT_TYPE,
		Event: LeaveChatEvent{
			UserId:  userId,
			ActorId: actorId,
		},
	}
	bus.messenger.Post(ew)
	bus.hub.LeaveChat(chatId, userId)
	log.Printf("Leave chat event sent to chat %d", chatId)
}

func (bus *EventBus) PostEnterChatEvent(chatId, userId, actorId int64) {
	ew := EventWrapper{
		ChatId:    chatId,
		EventType: ENTER_CHAT_EVENT_TYPE,
		Event: EnterChatEvent{
			UserId:  userId,
			ActorId: actorId,
		},
	}

	bus.hub.JoinChat(chatId, userId)
	bus.messenger.Post(ew)
	log.Printf("Join chat event sent to chat %d", chatId)
}

type InMemoryMessenger struct {
	channel chan EventWrapper
}

func NewInMemoryMessenger() *InMemoryMessenger {
	return &InMemoryMessenger{
		channel: make(chan EventWrapper, 64),
	}
}

func (m *InMemoryMessenger) Post(e EventWrapper) error {
	select {
	case m.channel <- e:
	default:
	}
	return nil
}

func (m *InMemoryMessenger) Listen(c context.Context, functions ...func(e EventWrapper) error) error {
	for {
		select {
		case event, ok := <-m.channel:
			if ok {
				for _, f := range functions {
					err := f(event)
					if err != nil {
						log.Print(err)
					}
				}
			}
		case <-c.Done():
			close(m.channel)
			return nil
		}
	}
}
