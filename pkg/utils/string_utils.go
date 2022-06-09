package utils

import (
	"git.internal.yunify.com/qxp/misc/logger"
	"strings"
)

// StringJoins string join
func StringJoins(strs ...string) string {
	var build strings.Builder
	for _, s := range strs {
		if len(s) > 0 {
			build.WriteString(s)
		}
	}
	return build.String()
}

// GetAsString 类型转换
func GetAsString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	logger.Logger.Errorf("invalid string value %#v", v)
	return ""
}
