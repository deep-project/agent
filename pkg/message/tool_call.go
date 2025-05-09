package message

import "encoding/json"

// ToolCall 工具调用
type ToolCall struct {
	ID        string            `json:"id,omitempty"`        // 工具调用ID
	ToolID    string            `json:"tool_id,omitempty"`   // 工具ID
	Arguments ToolCallArguments `json:"arguments,omitempty"` // 工具参数
}

type ToolCallArguments map[string]any

func (a *ToolCallArguments) String() string {
	if a == nil {
		return ""
	}
	if b, err := json.Marshal(a); err == nil {
		return string(b)
	}
	return ""
}

func (a *ToolCallArguments) Map() map[string]any {
	return *a
}

func NewToolCallArgumentsByString(val string) (args ToolCallArguments) {
	if err := json.Unmarshal([]byte(val), &args); err != nil {
		return nil
	}
	return
}
