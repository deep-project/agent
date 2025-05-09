package adapters

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/deep-project/agent/pkg/message"
	"github.com/deep-project/agent/pkg/mind"

	"github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	client    *openai.Client
	modelName string
}

func NewOpenAI(config openai.ClientConfig, modelName string) *OpenAI {
	return &OpenAI{
		client:    openai.NewClientWithConfig(config),
		modelName: modelName,
	}
}

func (o *OpenAI) Call(opt *mind.CallOptions) (*mind.CallResponse, error) {
	req := openai.ChatCompletionRequest{
		// 为什么不可以使用stream方式通信，经测试，stream方式同样会将tools的args也切割
		// 如果硬要处理，就要将tool角色单独摘出来拼凑之后统一返回，处理复杂度过大。
		//Stream: true,
		Model:    o.modelName,
		Tools:    o.convertToOpenAITools(opt.Tools),
		Messages: o.convertToOpenAIMessage(opt.Messages),
	}
	resp, err := o.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return nil, err
	}
	if len(resp.Choices) == 0 {
		return nil, errors.New("No response received")
	}
	choice := resp.Choices[0]
	return &mind.CallResponse{
		Message: *o.convertToAgentMessage(&choice.Message),
	}, nil
}

func (o *OpenAI) convertToOpenAITools(tools []mind.Tool) (res []openai.Tool) {
	for _, t := range tools {
		res = append(res, openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        t.ID,
				Description: t.Description,
				Parameters:  t.ParametersJSONSchema(),
			},
		})
	}
	return res
}

func (o *OpenAI) convertToAgentMessage(msg *openai.ChatCompletionMessage) (res *message.Message) {
	return &message.Message{
		Role:      message.Role(msg.Role),
		Contents:  o.convertToAgentMessageContent(msg),
		ToolCalls: o.convertToAgentToolCalls(&msg.ToolCalls),
	}
}

func (o *OpenAI) convertToAgentToolCalls(tools *[]openai.ToolCall) (res []message.ToolCall) {
	for _, t := range *tools {
		res = append(res, message.ToolCall{
			ID:        t.ID,
			ToolID:    t.Function.Name,
			Arguments: message.NewToolCallArgumentsByString(t.Function.Arguments),
		})
	}
	return res
}

func (o *OpenAI) convertToAgentMessageContent(msg *openai.ChatCompletionMessage) (res []message.Content) {
	if msg.Content != "" {
		res = append(res, message.NewMessageWithContentText(msg.Content))
	}
	for _, c := range msg.MultiContent {
		switch c.Type {
		case "text":
			res = append(res, message.NewMessageWithContentText(c.Text))
		case "image_url":
			res = append(res, message.NewMessageWithContentImage(c.ImageURL.URL))
		default:
			// 其他类型的处理
			// 直接格式化成文本返回吧，否则还能怎么办呢？
			if b, err := json.Marshal(&c); err == nil {
				res = append(res, message.NewMessageWithContentText(string(b)))
			}
		}
	}
	return res
}

func (o *OpenAI) convertToOpenAIMessage(msg []message.Message) (res []openai.ChatCompletionMessage) {
	for _, m := range msg {
		res = append(res, openai.ChatCompletionMessage{
			Role:         string(m.Role),
			MultiContent: o.convertToOpenAIMessageContent(m.Contents),
			ToolCalls:    o.convertToOpenAIToolCalls(m.ToolCalls),
			ToolCallID:   m.ToolCallID,
		})
	}
	return res
}

func (o *OpenAI) convertToOpenAIToolCalls(tools []message.ToolCall) (res []openai.ToolCall) {
	for _, t := range tools {
		res = append(res, openai.ToolCall{
			ID:   t.ID,
			Type: "function",
			Function: openai.FunctionCall{
				Name:      t.ToolID,
				Arguments: t.Arguments.String(),
			},
		})
	}
	return res
}

func (o *OpenAI) convertToOpenAIMessageContent(content []message.Content) (res []openai.ChatMessagePart) {
	for _, c := range content {
		switch c.Type {
		case message.ContentTypeText:
			res = append(res, openai.ChatMessagePart{Type: openai.ChatMessagePartTypeText, Text: c.Text.Text})
		case message.ContentTypeImage:
			res = append(res, openai.ChatMessagePart{Type: openai.ChatMessagePartTypeImageURL, ImageURL: &openai.ChatMessageImageURL{URL: c.Image.URI, Detail: openai.ImageURLDetail(c.Image.Detail)}})
		default:
			// 其他类型的处理
			// openai当前不支持其他类型
			// 直接格式化成文本返回吧，要相信大模型的智慧
			if b, err := json.Marshal(c); err == nil {
				res = append(res, openai.ChatMessagePart{Type: openai.ChatMessagePartTypeText, Text: string(b)})
			}
		}
	}
	return res
}
