package pubsub

import (
	"sync"

	pb "github.com/LuukBlankenstijn/fogistration/internal/shared/pb"
)

type Manager struct {
	mu          sync.RWMutex
	subscribers map[string]chan<- *pb.ServerMessage
}

func NewManager() *Manager {
	return &Manager{
		subscribers: make(map[string]chan<- *pb.ServerMessage),
	}
}

// Subscribe adds a subscriber for the given IP
func (m *Manager) Subscribe(ip string) <-chan *pb.ServerMessage {
	ch := make(chan *pb.ServerMessage, 16)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.subscribers[ip] = ch
	return ch
}

// Unsubscribe removes a subscriber for the given IP
func (m *Manager) Unsubscribe(ip string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if ch, exists := m.subscribers[ip]; exists {
		close(ch)
		delete(m.subscribers, ip)
	}
}

// Publish sends a message to a specific IP
func (m *Manager) Publish(ip string, msg *pb.ServerMessage) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if ch, exists := m.subscribers[ip]; exists {
		select {
		case ch <- msg:
		default:
			// Channel full, skip message
		}
	}
}

// Broadcast sends a message to all subscribers
func (m *Manager) Broadcast(msg *pb.ServerMessage) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, ch := range m.subscribers {
		select {
		case ch <- msg:
		default:
			// Channel full, skip message
		}
	}
}
