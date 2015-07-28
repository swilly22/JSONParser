package json_parser

import (
	"strings"
)

// stateFunc represents the state of the scanner
// as a function that returns the next state.
type stateFunc func(*lexer) stateFunc

func initState(l *lexer) stateFunc {
	r := l.peek()
	if r == leftMeta {
		return lexLeftMeta
	}

	return l.errorf("json must begin with '{'")
}

func lextText(l *lexer) stateFunc {
	for {
		if strings.HasPrefix(l.input[l.pos:], "{") {
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexLeftMeta // Next state.
		}
		if l.next() == eof {
			break
		}
	}
	// Correctly reached EOF.
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.emit(itemEOF)
	return nil
}

func lexLeftMeta(l *lexer) stateFunc {
	l.pos += len("{")
	l.emit(itemLeftMeta)
	return lexInsideObject // Now inside {}.
}

func lexRightMeta(l *lexer) stateFunc {
	l.pos += len("}")
	l.emit(itemRightMeta)
	return lexOutsideObject
}

func lexInsideObject(l *lexer) stateFunc {
	// Expecting either an identifier or '}'
	switch r := l.next(); {
	case r == eof:
		return l.errorf("unclosed object")
	case isAlphaNumeric(r):
		return lexIdentifier
	case r == rightMeta:
		return lexRightMeta
	}

	return l.errorf("invalid elemnt after left meta")
}

func lexOutsideObject(l *lexer) stateFunc {
	// Expecting either ',' or eof
	switch r := l.next(); {
	case r == eof:
		// we're done
		return nil
	case r == comma:
		return lexComma
	}
	return l.errorf("invalid elemnt after right meta")
}

func lexLeftBracket(l *lexer) stateFunc {
	l.pos += len("[")
	l.emit(itemLeftBracket)
	return lexInsideArray
}

func lexRightBracket(l *lexer) stateFunc {
	l.pos += len("]")
	l.emit(itemRightBracket)
	return lexOutsideArray
}

func lexInsideArray(l *lexer) stateFunc {
	// expecting either Number, String, '[', ']', '{'
	switch r := l.peek(); {
	case isNumeric(r):
		return lexNumber
	case r == quotationMark:
		return lexString
	case r == leftMeta:
		return lexLeftMeta
	case r == rightBracket:
		return lexRightBracket
	case r == leftBracket:
		return lexLeftBracket
	}

	return l.errorf("invalid element inside array")
}

func lexOutsideArray(l *lexer) stateFunc {
	// expecting either ',' or '}'
	switch r := l.next(); {
	case r == comma:
		return lexComma
	case r == rightMeta:
		return lexRightMeta
	}

	return l.errorf("invalid element after right bracket")
}

func lexComma(l *lexer) stateFunc {
	l.pos += len(",")
	l.emit(itemComma)

	// either Number, String, Identifier, '[' or '{'
	switch r := l.next(); {
	case isNumeric(r):
		l.backup()
		return lexNumber
	case r == quotationMark:
		l.backup()
		return lexString
	case isAlphaNumeric(r):
		return lexIdentifier
	case r == leftBracket:
		l.backup()
		return lexLeftBracket
	case r == leftMeta:
		l.backup()
		return lexLeftMeta
	}
	return l.errorf("invalid element after comma")
}

func lexNumber(l *lexer) stateFunc {
	// Optional leading sign.
	l.accept("+-")
	// Is it hex?
	digits := "0123456789"
	if l.accept("0") || l.accept("xX") {
		digits = "0123456789ABCDEF"
	}
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	l.emit(itemNumber)

	// expecting either ',' ']', '}'
	switch r := l.next(); {
	case r == comma:
		return lexComma
	case r == rightBracket:
		return lexRightBracket
	case r == rightMeta:
		return lexRightMeta
	}

	return l.errorf("ivalid element after number")
}

func lexString(l *lexer) stateFunc {

	// string opening quotation mark
	l.accept("\"")

	// consume string
	// TODO deal with escaping. (\")
	for strings.IndexRune("\"", l.next()) < 0 {
	}

	// string closing quotation mark
	l.accept("\"")

	l.emit(itemString)

	// Expecting either ',' , ']' or '}'
	switch r := l.next(); {
	case r == comma:
		return lexComma
	case r == rightBracket:
		return lexRightBracket
	case r == rightMeta:
		return lexRightMeta
	}

	return l.errorf("invalid element after string")
}

func lexIdentifier(l *lexer) stateFunc {
	// identifier must begin with an alphnumeric character
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if !l.accept(characters) {
		return l.errorf("identifier must begin with a character")
	}

	// consume identifier
	alphanumeric := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	l.acceptRun(alphanumeric)

	// identifier must end with ':'
	if !l.accept(":") {
		return l.errorf("identifier must end with ':'")
	}

	// expecting either String, Number '[' or '{'
	switch r := l.next(); {
	case isNumeric(r):
		l.backup()
		return lexNumber
	case r == quotationMark:
		l.backup()
		return lexString
	case r == leftBracket:
		return lexLeftBracket
	case r == leftMeta:
		return lexLeftMeta
	}

	return l.errorf("incalud element after identifier")
}
