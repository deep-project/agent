package ability

type Item struct {
	Name        string // name
	Description string // 描述
	Enable      bool   // 启用
	tools       []Tool // 工具列表，初始化item即初始化，可以减轻运行时初始化的性能消耗
	handler     Handler
}

func newItem(handler Handler) (*Item, error) {
	if handler == nil {
		return nil, ErrAbilityHandlerNotDefined
	}
	tools, err := handler.Tools()
	if err != nil {
		return nil, err
	}
	return &Item{
		Name:        handler.Name(),
		Description: handler.Description(),
		Enable:      handler.Enable(),
		tools:       tools,
		handler:     handler,
	}, nil
}

func (i *Item) Tools() []Tool {
	return i.tools
}

// 初始化tools,如果外部接口tools有更新，可以重新初始化
func (i *Item) InitTools() (err error) {
	i.tools, err = i.handler.Tools()
	return
}
