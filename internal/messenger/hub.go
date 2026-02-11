package messenger

import (
	"sync"

	"github.com/google/uuid"
)

type Connection struct {
	UserID int64
	Events chan EventWrapper
}

type ConnectionHub struct {
	//connection id -> channel event
	conns map[string]*Connection
	//chatId -> set of connection ids
	chatConns map[int64]map[string]struct{}

	mu *sync.RWMutex
}

func NewConnectionHub() *ConnectionHub {
	return &ConnectionHub{
		//connection id -> connection
		conns: make(map[string]*Connection),
		//chat id -> set of connection ids
		chatConns: make(map[int64]map[string]struct{}),

		mu: &sync.RWMutex{},
	}
}

func (h *ConnectionHub) CreateChatIfNotExist(chatID int64) {
	if _, ok := h.chatConns[chatID]; !ok {
		h.chatConns[chatID] = make(map[string]struct{})
	}
}

func (h *ConnectionHub) Subscribe(c chan EventWrapper, userID int64, chatIDs []int64) string {
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

func (h *ConnectionHub) LeaveChat(chatID, userID int64) {
	h.mu.Lock()
	for connID, connection := range h.conns {
		if connection.UserID == userID {
			connectionSet, ok := h.chatConns[chatID]
			if _, ok2 := connectionSet[connID]; ok && ok2 {
				delete(connectionSet, connID)
			}
		}
	}
	h.mu.Unlock()
}

func (h *ConnectionHub) JoinChat(chatID, userID int64) {
	h.mu.Lock()
	for connID, connection := range h.conns {
		if connection.UserID == userID {
			h.CreateChatIfNotExist(chatID)
			h.chatConns[chatID][connID] = struct{}{}
		}
	}
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

// This function should not be used by the owner of the websocket connection
//
// The projects architecture insists that CRUD operations and event emitting
// should be done by HTTP methods ONLY, not by websockets..
func (h *ConnectionHub) Broadcast(event EventWrapper) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	set, ok := h.chatConns[event.Chat.ID]
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
