package main

import (
	"fmt"
)

// item represents a token returned from the scanner.
type item struct {
	typ itemType // Type, such as itemNumber.
	val string   // Value, such as "23.2"
}

func (i item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return i.val
	}
	if len(i.val) > 10 {
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

// itemType identifies the type of lex items.
type itemType int

const (
	itemError itemType = iota // error occurred
	itemText
	itemLeftMeta
	itemRightMeta
	itemLeftBracket
	itemRightBracket
	itemIdentifier
	itemNumber
	itemString
	itemComma
	itemEOF
)

const (
	eof           = -1
	leftMeta      = '{'
	rightMeta     = '}'
	leftBracket   = '['
	rightBracket  = ']'
	comma         = ','
	quotationMark = '"'
)
