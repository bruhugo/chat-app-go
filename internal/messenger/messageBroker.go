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
		Chat:      *m.Chat,
		EventType: CREATE_MESSAGE_EVENT_TYPE,
		Event: CreateMessageEvent{
			Content:   m.Content,
			User:      *m.User,
			CreatedAt: m.CreatedAt,
		},
	}
	bus.messenger.Post(ew)
	log.Printf("Create message event of message %d sent to chat %d", m.ID, m.Chat.ID)
}

func (bus *EventBus) PostDeleteMessageEvent(messageId int64, chat models.Chat) {
	ew := EventWrapper{
		Chat:      chat,
		EventType: CREATE_MESSAGE_EVENT_TYPE,
		Event: DeleteMessageEvent{
			MessageId: messageId,
		},
	}

	bus.messenger.Post(ew)
	log.Printf("Delete message event of message %d sent to chat %d", messageId, chat.ID)
}

func (bus *EventBus) PostLeaveChatEvent(chat models.Chat, user, actor models.User) {
	ew := EventWrapper{
		Chat:      chat,
		EventType: LEAVE_CHAT_EVENT_TYPE,
		Event: LeaveChatEvent{
			User:  user,
			Actor: actor,
		},
	}
	bus.messenger.Post(ew)
	bus.hub.LeaveChat(chat.ID, user.ID)
	log.Printf("Leave chat event sent to chat %d", chat.ID)
}

func (bus *EventBus) PostEnterChatEvent(chat models.Chat, user, actor models.User) {
	ew := EventWrapper{
		Chat:      chat,
		EventType: ENTER_CHAT_EVENT_TYPE,
		Event: EnterChatEvent{
			User:  user,
			Actor: actor,
		},
	}

	bus.hub.JoinChat(chat.ID, user.ID)
	bus.messenger.Post(ew)
	log.Printf("Join chat event sent to chat %d", chat.ID)
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
