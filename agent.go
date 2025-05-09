package agent

import (
	"errors"
	"sync"

	"github.com/deep-project/agent/internal/helpers"
	"github.com/deep-project/agent/pkg/ability"
	"github.com/deep-project/agent/pkg/memory"
	"github.com/deep-project/agent/pkg/message"
	"github.com/deep-project/agent/pkg/mind"

	"github.com/google/uuid"
)

type Agent struct {
	mind    *mind.Mind       // 思维
	memory  *memory.Memory   // 记忆
	ability *ability.Ability // 能力
	mu      sync.Mutex
}

func New() *Agent {
	return &Agent{
		mind:    new(mind.Mind),
		memory:  new(memory.Memory),
		ability: new(ability.Ability),
	}
}

// GrantMind 给智能体赋予智慧
func (a *Agent) GrantMind(handler mind.Handler) *Agent {
	a.mind.SetHandler(handler)
	return a
}

// GrantMemory 给智能体赋予记忆
func (a *Agent) GrantMemory(handler memory.Handler) *Agent {
	a.memory.SetHandler(handler)
	return a
}

// GrantAbility 给智能体赋予能力
func (a *Agent) GrantAbility(handler ...ability.Handler) *Agent {
	return a.GrantAbilities(handler)
}

// GrantAbilities 赋予智能体能力
func (a *Agent) GrantAbilities(handlers []ability.Handler) *Agent {
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, handler := range handlers {
		a.ability.Add(handler)
	}
	return a
}

// ResetAbilities 重新赋予智能体能力
func (a *Agent) ResetAbilities(abilities []ability.Handler) *Agent {
	a.ClearAbilities()
	a.GrantAbilities(abilities)
	return a
}

// ClearAbilities 清除所有能力
func (a *Agent) ClearAbilities() *Agent {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.ability.Clear()
	return a
}

/////////

// ListMessages 获取消息列表
func (a *Agent) ListMessages(sessionID string, limit int) ([]message.Message, error) {
	return a.memory.ListMessages(sessionID, limit)
}

// AddMessages 添加消息
func (a *Agent) AddMessages(sessionID string, messages []message.Message) error {
	return a.memory.AddMessages(sessionID, messages)
}

// HasMessageSession 消息对话是否存在
func (a *Agent) HasMessageSession(sessionID string, messages []message.Message) (bool, error) {
	return a.memory.HasMessageSession(sessionID)
}

// Talk 只返回文本聊天信息
func (a *Agent) Talk(sessionID, text string) (sid string, res string, err error) {
	output, err := a.Send(sessionID, text)
	if err != nil {
		return
	}
	return output.SessionID, helpers.JoinTextMessageContents(output.Message.Contents), nil
}

// Send 返回完整消息结构体
func (a *Agent) Send(sessionID, text string) (*InteractOutput, error) {
	return a.Interact(&InteractInput{
		SessionID:    sessionID,
		MessageLimit: 50,
		Messages: []message.Message{
			{Role: message.RoleUser, Contents: []message.Content{message.NewMessageWithContentText(text)}},
		},
	})
}

// Interact 与agent交互
func (a *Agent) Interact(input *InteractInput) (output *InteractOutput, err error) {
	if input == nil {
		return nil, errors.New("interact input is empty")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if input.SessionID == "" {
		input.SessionID = uuid.New().String()
	}
	if err = a.AddMessages(input.SessionID, input.Messages); err != nil {
		return
	}
	callResponse, err := a.call(input)
	if err != nil {
		return
	}
	return &InteractOutput{
		SessionID: input.SessionID,
		Message:   callResponse.Message,
	}, nil
}

func (a *Agent) call(input *InteractInput) (_ *mind.CallResponse, err error) {
	messages, err := a.ListMessages(input.SessionID, input.MessageLimit)
	if err != nil {
		return
	}
	if len(messages) == 0 {
		return nil, errors.New("messages cannot be empty.")
	}
	tools, err := helpers.AbilityItemsToMindTools(a.ability.Items())
	if err != nil {
		return
	}
	resp, err := a.mind.Call(&mind.CallOptions{Messages: messages, Tools: tools})
	if err != nil {
		return
	}
	if resp == nil {
		return nil, errors.New("No response received")
	}
	if err = a.memory.AddMessage(input.SessionID, &resp.Message); err != nil {
		return
	}
	if len(resp.Message.ToolCalls) > 0 {
		for _, toolCall := range resp.Message.ToolCalls {
			toolCallMsg, err := a.execToolCall(&toolCall)
			if err != nil {
				continue
			}
			a.memory.AddMessage(input.SessionID, toolCallMsg)
		}
		return a.call(input)
	}
	return resp, nil
}

func (a *Agent) execToolCall(toolCall *message.ToolCall) (*message.Message, error) {
	itemIndex, toolName, err := helpers.ParseMindToolID(toolCall.ToolID)
	if err != nil {
		return nil, err
	}
	msg, err := a.ability.Call(itemIndex, toolName, &toolCall.Arguments)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, errors.New("response message is empty")
	}
	msg.ToolCallID = toolCall.ID
	return msg, nil
}

type InteractInput struct {
	SessionID    string            `json:"session_id"`
	Messages     []message.Message `json:"messages"`
	MessageLimit int               `json:"message_limit"` // 限制对话上文消息数
}

type InteractOutput struct {
	SessionID string          `json:"session_id"`
	Message   message.Message `json:"message"`
}
