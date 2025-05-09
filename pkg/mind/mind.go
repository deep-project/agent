package mind

import (
	"github.com/deep-project/agent/pkg/message"
	"github.com/deep-project/agent/pkg/tool"
)

type Handler interface {
	Call(opt *CallOptions) (*CallResponse, error)
}

type Mind struct {
	handler Handler
}

func (m *Mind) SetHandler(handler Handler) error {
	if handler == nil {
		return ErrMindHandlerNotDefined
	}
	m.handler = handler
	return nil
}

func (m *Mind) Call(opt *CallOptions) (*CallResponse, error) {
	if m.handler == nil {
		return nil, ErrMindHandlerNotDefined
	}
	return m.handler.Call(opt)
}

type CallOptions struct {
	Messages []message.Message
	Tools    []Tool
}

type CallResponse struct {
	Message message.Message
}

// Tool mind所需的tool结构需带唯一id
type Tool struct {
	ID string
	*tool.Tool
}
