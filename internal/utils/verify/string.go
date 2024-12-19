package verify

import (
	"strings"
)

// EmptyString 空字符串
func EmptyString(str string) bool {
	return strings.Trim(str, " ") == ""
}
