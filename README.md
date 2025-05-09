# 🦉 AI Agent

## 介绍 / Introduction

本项目使得创建一个AI智能体变得非常简单，能够赋予其思维（LLM）、记忆（Storage）和能力（Tools）。不仅可以接入当前流行的MCP服务，还可以自定义工具进行接入，极大提高了智能体创建的效率。

This project makes it very easy to create an AI agent, providing it with thinking (LLM), memory (Storage), and capabilities (Tools). It not only supports integration with popular MCP services but also allows for custom tool integrations, greatly enhancing the efficiency of creating an intelligent agent.

## 安装 / Installation

```go
go get github.com/deep-project/agent
```

## 使用 / Usage
```go

a := agent.New()

// 赋予AI思维
mindConfig := openai.DefaultConfig("TOKEN XXXXXXXXXXXXXXX")
mindConfig.BaseURL = os.Getenv("https://api.openai.com/v1")
a.GrantMind(adapters.NewOpenAI(mindConfig, "gpt-4"))

// 赋予AI记忆
a.GrantMemory(adapters.NewMemorySimpleAdapter(999))

// 赋予AI能力
mcpClient, _ := client.NewStdioMCPClient(`./mcp`, []string{}, "-y")
adapters.MCPAdapterInitializeClient(mcpClient)
a.GrantAbility(adapters.NewMCPAdapter(&adapters.MCPAdapterOptions{Enable: true}, mcpClient))

// 设置一个会话id
var sessionID = ""
sessionID, msg, _ := a.Talk(sessionID, "hello world!")

// 得到回复
fmt.Println("reply", msg)
```

## 更多 / More
#### 可以赋予多种能力 / Can empower with multiple capabilities
```go

// MCP 1
a.GrantAbility(adapters.NewMCPAdapter(&adapters.MCPAdapterOptions{Enable: true}, MCP_1))
// MCP 2
a.GrantAbility(adapters.NewMCPAdapter(&adapters.MCPAdapterOptions{Enable: true}, MCP_2))

```
> 能力不止于mcp服务，可以是任何符合程序接口的tools。可以直接自定义一个满足接口的结构体，整体打包，这样也不必再开启一个mcp服务了。

#### 内置的存储适配器 / Built-in storage adapter
```go
// 简单的存储(依靠内存)
a.GrantMemory(adapters.NewMemorySimpleAdapter(999))

// 使用bbolt数据库作为存储
db, _ := bbolt.Open("my.db", 0666, nil)
a.GrantMemory(adapters.NewMemoryBoltDBAdapter(db))
```
> 其他自定义的存储mysql sqlite pgsql都可以。只需要符合程序接口即可。

#### 多种交互方式 / Multiple interaction methods
```go
// 简单的文字交流
a.Talk(sessionID, "hello world!")

// 输入文字，返回完整消息体
a.Send(sessionID, "hello world!")

// 通过消息体交互
a.Interact(&agent.InteractInput{
		SessionID:    sessionID,
		MessagesLimit: 50,
		Messages: []message.Message{
			{Role: message.RoleUser, Contents: []message.Content{message.NewMessageWithContentText("hello world!")}},
		}
})

```
> 通过消息体交互，可以保持最大的灵活性，可以自定义角色，限制消息列表最大长度，发送多种类型的消息。