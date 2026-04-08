package ability

type Tool struct {
	Name        string
	Enable      bool // 启用
	Description string
	Parameters  []ToolParameter // 参数
	Required    []string        // 必填的字段集合
}

// Parameters Convert To JSON Schema
func (t *Tool) ParametersJSONSchema() *JSONSchema {
	res := &JSONSchema{
		Type:       "object",
		Properties: arr2Obj(t.Parameters),
		Required:   t.Required,
	}
	return res
}

func arr2Obj(arr []ToolParameter) map[string]ToolParameter {
	res := make(map[string]ToolParameter)
	for _, item := range arr {
		res[item.Name] = item
	}
	return res
}

// ToolParameterType 定义参数类型
type ToolParameterType string

// JSON Schema 支持的常用类型
const (
	ToolParameterTypeString  ToolParameterType = "string"
	ToolParameterTypeNumber  ToolParameterType = "number"
	ToolParameterTypeInteger ToolParameterType = "integer"
	ToolParameterTypeBoolean ToolParameterType = "boolean"
	ToolParameterTypeArray   ToolParameterType = "array"
	ToolParameterTypeObject  ToolParameterType = "object"
)

type ToolParameter struct {
	Name        string            `json:"name"`                  // 属性名
	Type        ToolParameterType `json:"type,omitempty"`        // 类型
	Description string            `json:"description,omitempty"` // 描述
	Title       string            `json:"title,omitempty"`       // 友好属性名称，可以显示更易读的属性名称
	Enum        []string          `json:"enum,omitempty"`        // 枚举值
	Default     any               `json:"default,omitempty"`     // 默认值
	MaxLength   int               `json:"maxLength,omitempty"`   // 属性值最大长度
	MinLength   int               `json:"minLength,omitempty"`   // 属性值最小长度
	Pattern     string            `json:"pattern,omitempty"`     // 属性值必须匹配正则表达式
	Maximum     float64           `json:"maximum,omitempty"`     // 属性值为数字时的最大值
	Minimum     float64           `json:"minimum,omitempty"`     // 属性值为数字时的最小值
	MultipleOf  float64           `json:"multipleOf,omitempty"`  // 属性值为数字时必须是指定倍数（数值必须能被此值整除）
}

type JSONSchema struct {
	Type       string   `json:"type"`
	Properties any      `json:"properties"`
	Required   []string `json:"required,omitempty"`
}
