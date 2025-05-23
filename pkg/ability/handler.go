package ability

import (
	"github.com/deep-project/agent/pkg/message"
)

type Handler interface {
	Name() string
	Description() string
	Enable() bool
	Tools() ([]Tool, error)
	CallTool(opt *CallToolOptions) (*message.Message, error)
}

type CallToolOptions struct {
	Name string
	Args *message.ToolCallArguments
	Meta Meta
}

type Meta map[string]any

func NewMeta() Meta {
	return make(map[string]any)
}
