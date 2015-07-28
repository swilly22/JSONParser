package main

import (
	"fmt"
)

// stateFunc represents the state of the scanner
// as a function that returns the next state.
type stateFunc func(*lexer) stateFunc

func initState(l *lexer) stateFunc {
	fmt.Println("in initState")

	r := l.peek()
	if r == leftMeta {
		return lexLeftMeta
	}

	return l.errorf("json must begin with '{'")
}

func lexLeftMeta(l *lexer) stateFunc {
	fmt.Println("in lexLeftMeta")

	l.pos += len("{")
	l.emit(itemLeftMeta)
	return lexInsideObject // Now inside {}.
}

func lexRightMeta(l *lexer) stateFunc {
	fmt.Println("in lexRightMeta")

	l.pos += len("}")
	l.emit(itemRightMeta)
	return lexOutsideObject
}

func lexInsideObject(l *lexer) stateFunc {
	fmt.Println("in lexInsideObject")

	// Expecting either an identifier or '}'
	switch r := l.peek(); {
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
	fmt.Println("in lexOutsideObject")

	// Expecting either ',', ']', '}' or eof
	switch r := l.peek(); {
	case r == eof:
		// we're done
		return nil
	case r == comma:
		return lexComma
	case r == rightBracket:
		return lexRightBracket
	case r == rightMeta:
		return lexRightMeta
	}
	return l.errorf("invalid elemnt after right meta")
}

func lexLeftBracket(l *lexer) stateFunc {
	fmt.Println("in lexLeftBracket")

	l.pos += len("[")
	l.emit(itemLeftBracket)
	return lexInsideArray
}

func lexRightBracket(l *lexer) stateFunc {
	fmt.Println("in lexRightBracket")

	l.pos += len("]")
	l.emit(itemRightBracket)
	return lexOutsideArray
}

func lexInsideArray(l *lexer) stateFunc {
	fmt.Println("in lexInsideArray")

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
	fmt.Println("in lexOutsideArray")

	// expecting either ',', ']' or '}'
	switch r := l.peek(); {
	case r == comma:
		return lexComma
	case r == rightMeta:
		return lexRightMeta
	case r == rightBracket:
		return lexRightBracket
	}

	return l.errorf("invalid element after right bracket")
}

func lexComma(l *lexer) stateFunc {
	fmt.Println("in lexComma")

	l.pos += len(",")
	l.emit(itemComma)

	// either Number, String, Identifier, '[' or '{'
	switch r := l.peek(); {
	case isNumeric(r):
		return lexNumber
	case r == quotationMark:
		return lexString
	case isAlphaNumeric(r):
		return lexIdentifier
	case r == leftBracket:
		return lexLeftBracket
	case r == leftMeta:
		return lexLeftMeta
	}
	return l.errorf("invalid element after comma")
}

func lexNumber(l *lexer) stateFunc {
	fmt.Println("in lexNumber")

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
	switch r := l.peek(); {
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
	fmt.Println("in lexString")

	// string opening quotation mark
	if !l.accept("\"") {
		l.errorf("string should be enclosed with quotation marks")
	}

	// consume string
	// TODO deal with escaping. (\")
	alphanumeric := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	l.acceptRun(alphanumeric)

	// string closing quotation mark
	if !l.accept("\"") {
		l.errorf("string should be enclosed with quotation marks")
	}

	l.emit(itemString)

	// Expecting either ',' , ']' or '}'
	switch r := l.peek(); {
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
	fmt.Println("in lexIdentifier")

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

	l.emit(itemIdentifier)

	// expecting either String, Number '[' or '{'
	switch r := l.peek(); {
	case isNumeric(r):
		return lexNumber
	case r == quotationMark:
		return lexString
	case r == leftBracket:
		return lexLeftBracket
	case r == leftMeta:
		return lexLeftMeta
	}

	return l.errorf("incalud element after identifier")
}
