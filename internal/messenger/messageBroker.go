package messenger

import (
	"context"
)

type Messenger interface {
	Post(Event) error
	Subscribe(chan Event) error
	Unsubscribe() error
}

type EventBus struct {
	hub       *ConnectionHub
	messenger Messenger
	channel   chan Event
}

func NewEventBus(m Messenger, h *ConnectionHub, c context.Context) *EventBus {
	ma := &EventBus{
		messenger: m,
		hub:       h,
	}

	return ma
}

func (eb *EventBus) Post(e Event) {
	eb.messenger.Post(e)
}

type InMemoryMessenger struct {
	channel chan Event
}

func NewInMemoryMessenger() *InMemoryMessenger {
	return &InMemoryMessenger{
		channel: make(chan Event),
	}
}

func (m *InMemoryMessenger) Post(e Event) {
	m.channel <- e
}

func (m *InMemoryMessenger) Listen(c context.Context, functions ...func(e Event) error) {
	for {
		select {
		case event, ok := <-m.channel:
			if ok {
				for _, f := range functions {
					f(event)
				}
			}
		case <-c.Done():
			close(m.channel)
			return
		}
	}
}
