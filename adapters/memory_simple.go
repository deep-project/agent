package adapters

import (
	"sync"

	"github.com/deep-project/agent/pkg/message"
)

type MemorySimpleAdapter struct {
	MaxSize int

	store map[string][]message.Message
	mu    sync.RWMutex
}

func NewMemorySimpleAdapter(maxSize int) *MemorySimpleAdapter {
	return &MemorySimpleAdapter{
		MaxSize: maxSize,
		store:   make(map[string][]message.Message),
	}
}

func (m *MemorySimpleAdapter) HasMessageSession(sessionID string) (bool, error) {
	_, exists := m.store[sessionID]
	return exists, nil
}

func (m *MemorySimpleAdapter) AddMessage(sessionID string, msg *message.Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	messages := m.store[sessionID]
	if messages != nil && m.MaxSize > 0 && len(messages) >= m.MaxSize {
		messages = messages[1:] // 如果消息数量超过最大限制，则删除最早的一条消息
	}
	m.store[sessionID] = append(messages, *msg)
	return nil
}

func (m *MemorySimpleAdapter) ListMessages(sessionID string, limit int) ([]message.Message, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	list, exist := m.store[sessionID]
	if !exist {
		return []message.Message{}, nil
	}
	if limit > 0 && len(list) > limit {
		return list[:limit], nil
	}
	return list, nil
}
