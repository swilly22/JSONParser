package json_parser

import (
	"strings"
)

func isNumeric(r rune) bool {
	digits := "+-0987654321"
	return strings.IndexRune(digits, r) >= 0
}

func isAlphaNumeric(r rune) bool {
	digits := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return strings.IndexRune(digits, r) >= 0
}
