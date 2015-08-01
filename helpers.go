package main

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

func matchExact(valid string, l *lexer) bool {
	if len(valid) > len(l.input[l.pos:]) {
		return false
	}

	backupCount := 0
	idx := 0

	for ; idx < len(valid); idx++ {
		backupCount++
		if strings.IndexRune(valid[idx:idx+1], l.next()) == -1 {
			break
		}
	}

	// backup
	for ; backupCount > 0; backupCount-- {
		l.backup()
	}

	return idx == len(valid)
}

func isTrue(l *lexer) bool {
	return matchExact("true", l)
}

func isFalse(l *lexer) bool {
	return matchExact("false", l)
}

func isNull(l *lexer) bool {
	return matchExact("null", l)
}
