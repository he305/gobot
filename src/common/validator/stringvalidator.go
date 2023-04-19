package validator

import "strings"

func ValidateString(str string) bool {
	return len(strings.TrimSpace(str)) > 0
}
