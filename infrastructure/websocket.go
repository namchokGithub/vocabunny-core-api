package infrastructure

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WebsocketManager struct {
	mu    sync.RWMutex
	conns map[string]*websocket.Conn
}

func NewWebsocketManager() *WebsocketManager {
	return &WebsocketManager{
		conns: make(map[string]*websocket.Conn),
	}
}

func (m *WebsocketManager) Add(key string, conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.conns[key] = conn
}

func (m *WebsocketManager) Remove(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.conns, key)
}

func (m *WebsocketManager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.conns)
}
