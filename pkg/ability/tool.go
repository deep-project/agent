package ability

type Tool struct {
	Name        string
	Enable      bool // 启用
	Description string
	Parameters  []ToolParameter // 参数
}

// Parameters Convert To JSON Schema
func (t *Tool) ParametersJSONSchema() *JSONSchema {
	res := &JSONSchema{
		Type:       "object",
		Properties: t.Parameters,
		Required:   []string{},
	}
	for _, p := range t.Parameters {
		if p.Required {
			res.Required = append(res.Required, p.Name)
		}
	}
	return res
}

type ToolParameter struct {
	Name        string   `json:"name"`                  // 属性名
	Type        string   `json:"type,omitempty"`        // 类型
	Description string   `json:"description,omitempty"` // 描述
	Title       string   `json:"title,omitempty"`       // 友好属性名称，可以显示更易读的属性名称
	Required    bool     `json:"required,omitempty"`    // 是否必填
	Enum        []string `json:"enum,omitempty"`        // 枚举值
	Default     any      `json:"default,omitempty"`     // 默认值
	MaxLength   int      `json:"maxLength,omitempty"`   // 属性值最大长度
	MinLength   int      `json:"minLength,omitempty"`   // 属性值最小长度
	Pattern     string   `json:"pattern,omitempty"`     // 属性值必须匹配正则表达式
	Maximum     float64  `json:"maximum,omitempty"`     // 属性值为数字时的最大值
	Minimum     float64  `json:"minimum,omitempty"`     // 属性值为数字时的最小值
	MultipleOf  float64  `json:"multipleOf,omitempty"`  // 属性值为数字时必须是指定倍数（数值必须能被此值整除）
}

type JSONSchema struct {
	Type       string   `json:"type"`
	Properties any      `json:"properties"`
	Required   []string `json:"required,omitempty"`
}
