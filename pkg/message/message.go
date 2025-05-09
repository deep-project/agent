package message

type Message struct {
	Role       Role       `json:"role"`                   // 消息角色
	Contents   []Content  `json:"content,omitempty"`      // 消息内容
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`   // 如果是assistant角色，可能有需要调用的工具列表
	ToolCallID string     `json:"tool_call_id,omitempty"` // 如果是tool角色，需设定ToolCallID
}
