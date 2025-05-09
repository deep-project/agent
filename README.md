# AI Agent

## 介绍 / Introduction

本项目可以非常方面的创建一个AI Agent，并且赋予其思维(LLM)，赋予其记忆(Storage)，赋予其能力(Tools)。


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