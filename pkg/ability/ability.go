package ability

import (
	"sync"

	"github.com/deep-project/agent/pkg/message"
)

type Ability struct {
	items []Item
	mu    sync.RWMutex
}

func (a *Ability) Add(handler Handler) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	item, err := newItem(handler)
	if err != nil {
		return err
	}
	a.items = append(a.items, *item)
	return nil
}

func (a *Ability) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.items = []Item{}
}

func (a *Ability) Items() []Item {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.items
}

func (a *Ability) Call(index int, toolName string, args *message.ToolCallArguments, meta Meta) (_ *message.Message, err error) {
	item, err := a.getItem(index)
	if err != nil {
		return
	}
	if item.handler == nil {
		return nil, ErrAbilityHandlerNotDefined
	}
	return item.handler.CallTool(&CallToolOptions{Name: toolName, Args: args, Meta: meta})
}

func (a *Ability) getItem(index int) (*Item, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if index > len(a.items)-1 {
		return nil, ErrAbilityItemNotFound
	}
	return &a.items[index], nil
}
