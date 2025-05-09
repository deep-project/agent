package test

import (
	"os"

	"github.com/deep-project/agent"
	"github.com/deep-project/agent/adapters"
	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/client"
	"github.com/sashabaranov/go-openai"
	"go.etcd.io/bbolt"
)

func useMind(a *agent.Agent) {
	godotenv.Load(".env")

	mindConfig := openai.DefaultConfig(os.Getenv("TOKEN"))
	mindConfig.BaseURL = os.Getenv("BASE_URL")
	a.GrantMind(adapters.NewOpenAI(mindConfig, os.Getenv("MODEL")))
}

func useSimpleMemory(a *agent.Agent) {
	a.GrantMemory(adapters.NewMemorySimpleAdapter(999))
}

func useBoltDBMemory(a *agent.Agent) (err error) {
	db, err := bbolt.Open("test.db", 0666, nil)
	if err != nil {
		return
	}
	a.GrantMemory(adapters.NewMemoryBoltDBAdapter(db))
	return nil
}

func useAbility(a *agent.Agent) (err error) {
	mcpClient, err := client.NewStdioMCPClient(`./mcp`, []string{}, "-y")
	if err != nil {
		return
	}
	if _, err = adapters.MCPAdapterInitializeClient(mcpClient); err != nil {
		return
	}
	a.GrantAbility(adapters.NewMCPAdapter(&adapters.MCPAdapterOptions{Enable: true}, mcpClient))
	return
}
