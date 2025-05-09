package message

// 消息内容
type Content struct {
	Type     ContentType     `json:"type"`
	Text     ContentText     `json:"text,omitempty"`
	Image    ContentImage    `json:"image,omitempty"`
	Resource ContentResource `json:"resource,omitempty"`
}

type ContentType string

const (
	ContentTypeText     ContentType = "text"
	ContentTypeImage    ContentType = "image"
	ContentTypeResource ContentType = "resource"
)

type ContentText struct {
	Text string `json:"text,omitempty"`
}

type ContentImage struct {
	URI    string `json:"uri,omitempty"`    // 可以是url或base64编码
	Detail string `json:"detail,omitempty"` // 图片质量 high low auto
}

type ContentResource struct {
	MIMEType string `json:"mimeType"` // 资源类型
	Data     []byte `json:"data"`     // 资源数据
}

func NewMessageWithContentText(text string) Content {
	return Content{Type: ContentTypeText, Text: ContentText{Text: text}}
}

func NewMessageWithContentImage(uri string) Content {
	return Content{Type: ContentTypeImage, Image: ContentImage{URI: uri, Detail: "auto"}}
}
