package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/deep-project/agent/pkg/message"
	"github.com/deep-project/agent/pkg/tool"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

type MCPAdapterOptions struct {
	Name        string
	Description string
	Enable      bool
	Timeout     time.Duration
}

type MCPAdapter struct {
	options *MCPAdapterOptions
	client  client.MCPClient
}

func NewMCPAdapter(options *MCPAdapterOptions, cli client.MCPClient) *MCPAdapter {
	if options.Timeout == 0 {
		options.Timeout = 30 * time.Second
	}
	return &MCPAdapter{options: options, client: cli}
}

func MCPAdapterInitializeClient(cli client.MCPClient) (*mcp.InitializeResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	initReq := mcp.InitializeRequest{}
	initReq.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initReq.Params.ClientInfo = mcp.Implementation{Name: "deep-project/agent", Version: "1.0.0"}
	return cli.Initialize(ctx, initReq)
}

func (m *MCPAdapter) Name() string {
	return m.options.Name
}

func (m *MCPAdapter) Enable() bool {
	return m.options.Enable
}

func (m *MCPAdapter) Description() string {
	return m.options.Description
}

func (m *MCPAdapter) Tools() (res []tool.Tool, _ error) {
	toolsRequest := mcp.ListToolsRequest{}
	list, err := m.client.ListTools(context.Background(), toolsRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to query list from MCP adapter, %s", err.Error())
	}
	for _, item := range list.Tools {
		_tool, err := m.convertToAgentTool(&item)
		if err != nil {
			//TODO
			// 如果转换工具出错，是跳过呢，还是直接返回错误呢？
			// 如果直接返回错误，会影响其他工具转换，毕竟外部mcp是不可控的，
			// 是否需要忍受这种个别错误？
			continue
		}
		res = append(res, *_tool)
	}
	return
}

func (m *MCPAdapter) CallTool(name string, args *message.ToolCallArguments) (*message.Message, error) {
	listDirRequest := mcp.CallToolRequest{Request: mcp.Request{Method: "tools/call"}}
	listDirRequest.Params.Name = name
	listDirRequest.Params.Arguments = args.Map()

	ctx, cancel := context.WithTimeout(context.Background(), m.options.Timeout)
	defer cancel() // 确保退出前释放资源
	result, err := m.client.CallTool(ctx, listDirRequest)
	if err != nil {
		return nil, err
	}
	contents, err := m.convertToAgentMessageContents(result.Content)
	if err != nil {
		return nil, err
	}
	return &message.Message{
		Role:     message.RoleTool,
		Contents: contents,
	}, nil
}

func (m *MCPAdapter) convertToAgentTool(mTool *mcp.Tool) (*tool.Tool, error) {
	parameters, err := m.convertToAgentToolParameters(mTool.InputSchema)
	if err != nil {
		return nil, err
	}
	return &tool.Tool{
		Name:        mTool.Name,
		Description: mTool.Description,
		Enable:      true,
		Parameters:  parameters,
	}, nil
}

func (m *MCPAdapter) convertToAgentToolParameters(inputSchema mcp.ToolInputSchema) (res []tool.Parameter, err error) {
	if inputSchema.Type != "object" || inputSchema.Properties == nil {
		return nil, errors.New("The input schema is malformed, or no properties are defined.")
	}

	requiredSet := make(map[string]bool)
	for _, name := range inputSchema.Required {
		requiredSet[name] = true
	}

	for name, prop := range inputSchema.Properties {
		param := tool.Parameter{
			Name:     name,
			Required: requiredSet[name],
		}

		// prop 是 map[string]interface{}
		propMap, ok := prop.(map[string]interface{})
		if !ok {
			continue // 如果结构不对就跳过
		}

		// 映射通用字段
		if v, ok := propMap["type"].(string); ok {
			param.Type = v
		}
		if v, ok := propMap["title"].(string); ok {
			param.Title = v
		}
		if v, ok := propMap["description"].(string); ok {
			param.Description = v
		}
		if v, ok := propMap["default"]; ok {
			param.Default = v
		}
		if v, ok := propMap["enum"].([]interface{}); ok {
			for _, e := range v {
				if s, ok := e.(string); ok {
					param.Enum = append(param.Enum, s)
				}
			}
		}
		if v, ok := propMap["maxLength"].(float64); ok {
			param.MaxLength = int(v)
		}
		if v, ok := propMap["minLength"].(float64); ok {
			param.MinLength = int(v)
		}
		if v, ok := propMap["pattern"].(string); ok {
			param.Pattern = v
		}
		if v, ok := propMap["maximum"].(float64); ok {
			param.Maximum = v
		}
		if v, ok := propMap["minimum"].(float64); ok {
			param.Minimum = v
		}
		if v, ok := propMap["multipleOf"].(float64); ok {
			param.MultipleOf = v
		}
		res = append(res, param)
	}
	return
}

func (m *MCPAdapter) convertToAgentMessageContents(mContents []mcp.Content) (res []message.Content, err error) {

	type typeHolder struct {
		Type string `json:"type"`
	}

	for _, mContent := range mContents {
		var holder typeHolder
		// 先将 content 转成 JSON 字节
		raw, err := json.Marshal(mContent)
		if err != nil {
			return nil, fmt.Errorf("marshal content: %w", err)
		}
		// 解析类型
		if err := json.Unmarshal(raw, &holder); err != nil {
			return nil, fmt.Errorf("unmarshal type field: %w", err)
		}
		switch holder.Type {
		case "text":
			var txt mcp.TextContent
			if err := json.Unmarshal(raw, &txt); err != nil {
				return nil, fmt.Errorf("unmarshal text content: %w", err)
			}
			res = append(res, message.NewMessageWithContentText(txt.Text))
		case "image":
			var img mcp.ImageContent
			if err := json.Unmarshal(raw, &img); err != nil {
				return nil, fmt.Errorf("unmarshal image content: %w", err)
			}
			res = append(res, message.NewMessageWithContentImage(img.Data))
		default:
			// TODO 其他类型通过text描述返回
			// 当前使用mcp库，对于如何处理resource类型的解析仍需推敲
			res = append(res, message.NewMessageWithContentText(string(raw)))
		}
	}
	return
}
