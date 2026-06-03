package utils

import (
	"tmeet/internal/exception"
	"unicode/utf8"
)

// CharacterLimit character limit
func CharacterLimit(flag, str string, limit int) error {
	if n := utf8.RuneCountInString(str); n > limit {
		return exception.InvalidArgsError.With("flag:%s, character limit exceeded, limit: %d, actual: %d", flag, limit, n)
	}
	return nil
}
