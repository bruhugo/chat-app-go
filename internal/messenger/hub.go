package messenger

import (
	"sync"

	"github.com/google/uuid"
)

type Connection struct {
	UserID string
	Events chan Event
}

type ConnectionHub struct {
	//connection id -> channel event
	conns map[string]*Connection
	//chatId -> set of connection ids
	chatConns map[string]map[string]struct{}

	mu *sync.RWMutex
}

func NewConnectionHub() *ConnectionHub {
	return &ConnectionHub{
		//connection id -> connection
		conns: make(map[string]*Connection),
		//chat id -> set of connection ids
		chatConns: make(map[string]map[string]struct{}),

		mu: &sync.RWMutex{},
	}
}

func (h *ConnectionHub) CreateChatIfNotExist(chatID string) {
	if _, ok := h.chatConns[chatID]; !ok {
		h.chatConns[chatID] = make(map[string]struct{})
	}
}

func (h *ConnectionHub) Subscribe(c chan Event, userID string, chatIDs []string) string {
	connectionID := uuid.NewString()
	h.mu.Lock()
	h.conns[connectionID] = &Connection{
		UserID: userID,
		Events: c,
	}

	for _, chatID := range chatIDs {
		h.CreateChatIfNotExist(chatID)
		h.chatConns[chatID][connectionID] = struct{}{}
	}
	h.mu.Unlock()

	return connectionID
}

func (h *ConnectionHub) LeaveChat(connectionID, chatID string) {
	h.mu.Lock()
	delete(h.chatConns[chatID], connectionID)

	m, ok := h.chatConns[chatID]
	if ok && len(m) == 0 {
		delete(h.chatConns, chatID)
	}

	h.mu.Unlock()
}

func (h *ConnectionHub) JoinChat(ConnectionID, chatID string) {
	h.mu.Lock()
	h.CreateChatIfNotExist(chatID)
	h.chatConns[chatID][ConnectionID] = struct{}{}
	h.mu.Unlock()
}

func (h *ConnectionHub) Unsubscribe(connectionID string) {
	h.mu.Lock()
	connection, ok := h.conns[connectionID]
	if !ok {
		return
	}
	close(connection.Events)
	delete(h.conns, connectionID)

	delete(h.conns, connectionID)
	for n, v := range h.chatConns {
		delete(v, connectionID)
		if len(v) == 0 {
			delete(h.chatConns, n)
		}
	}

	h.mu.Unlock()
}

func (h *ConnectionHub) Broadcast(event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	set, ok := h.chatConns[event.ChatId()]
	if !ok {
		return
	}

	for connectionId := range set {
		conn, ok := h.conns[connectionId]
		if !ok {
			continue
		}

		select {
		case conn.Events <- event:
		default:
		}
	}
}
