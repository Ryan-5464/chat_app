package credentials

import "regexp"

func NewUsername(s string) bool {
	if len(s) > 20 {
		return false
	}
	matched, _ := regexp.MatchString(`^[A-Za-z0-9_]+$`, s)
	return matched
}
