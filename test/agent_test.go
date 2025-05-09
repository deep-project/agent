package test

import (
	"testing"

	"github.com/deep-project/agent"
)

func TestAgent(t *testing.T) {
	a := agent.New()

	useMind(a)
	//useSimpleMemory(a)
	useBoltDBMemory(a)
	err := useAbility(a)
	if err != nil {
		t.Error(err)
		return
	}

	var sessionID = ""
	sessionID, msg, err := a.Talk(sessionID, "你好，帮我查一下180154有货吗？")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(sessionID, msg)
}
