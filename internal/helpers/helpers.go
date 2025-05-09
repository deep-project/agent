package helpers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/deep-project/agent/pkg/ability"
	"github.com/deep-project/agent/pkg/message"
	"github.com/deep-project/agent/pkg/mind"
)

// JoinTextMessages 提取消息中的所有文本类型内容并拼接成一个字符串
func JoinTextMessageContents(contents []message.Content) string {
	var res string
	for _, c := range contents {
		if c.Type == message.ContentTypeText {
			res += c.Text.Text
		}
	}
	return res
}

// ability items 转换成 Mind Tools
func AbilityItemsToMindTools(items []ability.Item) (res []mind.Tool, err error) {
	for i, item := range items {
		if !item.Enable {
			continue
		}
		for _, tool := range item.Tools() {
			if tool.Enable {
				res = append(res, mind.Tool{ID: GenerateMindToolID(i, tool.Name), Tool: &tool})
			}
		}
	}
	return res, nil
}

// 使用ability index和tool name 生成 mind tool id
func GenerateMindToolID(abilityIndex int, toolName string) string {
	return fmt.Sprintf("%d-%s", abilityIndex, toolName)
}

// 将 mind tool id 还原回 ability items index 和 tool name
func ParseMindToolID(id string) (int, string, error) {
	parts := strings.SplitN(id, "-", 2) // 只分成两部分
	if len(parts) != 2 {
		return 0, "", errors.New("invalid mind tool ID format")
	}
	abilityItemsIndex, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", errors.New("invalid ability index")
	}
	return abilityItemsIndex, parts[1], nil
}
