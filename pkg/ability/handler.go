package ability

import (
	"github.com/deep-project/agent/pkg/message"
	"github.com/deep-project/agent/pkg/tool"
)

type Handler interface {
	Name() string
	Description() string
	Enable() bool
	Tools() ([]tool.Tool, error)
	CallTool(name string, args *message.ToolCallArguments) (*message.Message, error)
}
