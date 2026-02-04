package messenger

import "sync"

type Messenger interface {
	Subscribe(chan Event)
	Post(Event)
}

type InMemoryMessageBroker struct {
	subs []chan Event
	mu   *sync.RWMutex
}

func (mb *InMemoryMessageBroker) Subscribe(c chan Event) {
	mb.mu.Lock()
	mb.subs = append(mb.subs, c)
	mb.mu.Unlock()
}

func (mb *InMemoryMessageBroker) Post(e Event) {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	for _, c := range mb.subs {
		c <- e
	}
}
