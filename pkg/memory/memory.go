package memory

import (
	"github.com/deep-project/agent/pkg/message"
)

type Handler interface {
	AddMessage(sessionID string, msg *message.Message) error
	ListMessages(sessionID string, limit int) ([]message.Message, error)
	HasMessageSession(sessionID string) (bool, error) // 消息对话是否存在
}

type Memory struct {
	handler Handler
}

func (m *Memory) SetHandler(handler Handler) error {
	if handler == nil {
		return ErrMemoryHandlerNotDefined
	}
	m.handler = handler
	return nil
}

func (m *Memory) AddMessages(sessionID string, messages []message.Message) (err error) {
	for _, msg := range messages {
		if err = m.AddMessage(sessionID, &msg); err != nil {
			return
		}
	}
	return nil
}

func (m *Memory) AddMessage(sessionID string, msg *message.Message) error {
	if m.handler == nil {
		return ErrMemoryHandlerNotDefined
	}
	return m.handler.AddMessage(sessionID, msg)
}

func (m *Memory) ListMessages(sessionID string, limit int) ([]message.Message, error) {
	if m.handler == nil {
		return nil, ErrMemoryHandlerNotDefined
	}
	return m.handler.ListMessages(sessionID, limit)
}

func (m *Memory) HasMessageSession(sessionID string) (bool, error) {
	if m.handler == nil {
		return false, ErrMemoryHandlerNotDefined
	}
	return m.handler.HasMessageSession(sessionID)
}
