package main

import "testing"

func validateTokens(expectedToken, emittedTokens []itemType, t *testing.T) bool {
	if len(expectedToken) != len(emittedTokens) {
		t.Errorf("expecting %d token to be emitted, got %d\r\n",
			len(expectedToken), len(emittedTokens))
		return false
	}

	for i := 0; i < len(expectedToken); i++ {
		if expectedToken[i] != emittedTokens[i] {
			t.Errorf("token #%d should have been %v, got %v\r\n",
				i, expectedToken[i], emittedTokens[i])

			return false
		}
	}

	return true
}

func TestLexEmptyObject(t *testing.T) {
	expectedToken := []itemType{itemLeftMeta, itemRightMeta, itemEOF}
	emittedTokens := lexJSON("{}")
	validateTokens(expectedToken, emittedTokens, t)
}

func TestLexEmptyObjects(t *testing.T) {
	expectedToken := []itemType{itemLeftMeta,
		itemIdentifier, itemColon, itemLeftMeta,
		itemIdentifier, itemColon, itemLeftMeta,
		itemIdentifier, itemColon, itemLeftMeta,
		itemIdentifier, itemColon, itemLeftMeta,
		itemRightMeta, itemRightMeta, itemRightMeta, itemRightMeta, itemRightMeta, itemEOF}
	emittedTokens := lexJSON("{a:{b:{c:{d:{}}}}}")
	validateTokens(expectedToken, emittedTokens, t)
}

func TestLexEmptyArray(t *testing.T) {
	expectedToken := []itemType{itemLeftMeta, itemIdentifier, itemColon, itemLeftBracket, itemRightBracket, itemRightMeta, itemEOF}
	emittedTokens := lexJSON("{a:[]}")
	validateTokens(expectedToken, emittedTokens, t)
}

func TestLexEmptyArrays(t *testing.T) {
	expectedToken := []itemType{itemLeftMeta, itemIdentifier, itemColon,
		itemLeftBracket, itemLeftBracket, itemLeftBracket, itemLeftBracket, itemLeftBracket,
		itemRightBracket, itemRightBracket, itemRightBracket, itemRightBracket, itemRightBracket,
		itemRightMeta, itemEOF}
	emittedTokens := lexJSON("{a:[[[[[]]]]]}")
	validateTokens(expectedToken, emittedTokens, t)
}

func TestLexNumericIdentifier(t *testing.T) {
	expectedToken := []itemType{itemLeftMeta, itemIdentifier, itemColon, itemNumber, itemRightMeta, itemEOF}
	JSON := "{a:1}"
	emittedTokens := lexJSON(JSON)
	validateTokens(expectedToken, emittedTokens, t)
}

func TestLexStringIdentifier(t *testing.T) {
	expectedToken := []itemType{itemLeftMeta, itemIdentifier, itemColon, itemString, itemRightMeta, itemEOF}
	JSON := "{b:\"value\"}"
	emittedTokens := lexJSON(JSON)
	validateTokens(expectedToken, emittedTokens, t)
}

func TestLexTrueIdentifier(t *testing.T) {
	expectedToken := []itemType{itemLeftMeta, itemIdentifier, itemColon, itemTrue, itemRightMeta, itemEOF}
	JSON := "{a:true}"
	emittedTokens := lexJSON(JSON)
	validateTokens(expectedToken, emittedTokens, t)
}

func TestLexFalseIdentifier(t *testing.T) {
	expectedToken := []itemType{itemLeftMeta, itemIdentifier, itemColon, itemFalse, itemRightMeta, itemEOF}
	JSON := "{a:false}"
	emittedTokens := lexJSON(JSON)
	validateTokens(expectedToken, emittedTokens, t)
}

func TestLexNullIdentifier(t *testing.T) {
	expectedToken := []itemType{itemLeftMeta, itemIdentifier, itemColon, itemNull, itemRightMeta, itemEOF}
	JSON := "{a:null}"
	emittedTokens := lexJSON(JSON)
	validateTokens(expectedToken, emittedTokens, t)
}

func TestLexTwoIdentifiers(t *testing.T) {
	expectedToken := []itemType{itemLeftMeta, itemIdentifier, itemColon, itemNumber, itemComma, itemIdentifier, itemColon, itemString, itemRightMeta, itemEOF}
	JSON := "{a:1,b:\"value\"}"
	emittedTokens := lexJSON(JSON)
	validateTokens(expectedToken, emittedTokens, t)
}

func TestLexArrayIdentifier(t *testing.T) {
	expectedToken := []itemType{itemLeftMeta, itemIdentifier, itemColon, itemLeftBracket,
		itemNumber, itemComma,
		itemTrue, itemComma,
		itemFalse, itemComma,
		itemNull,
		itemRightBracket, itemRightMeta, itemEOF}
	JSON := "{a:[1,true,false,null]}"
	emittedTokens := lexJSON(JSON)
	validateTokens(expectedToken, emittedTokens, t)
}

func TestLexArrayOfObjects(t *testing.T) {
	expectedToken := []itemType{itemLeftMeta, itemIdentifier, itemColon, itemLeftBracket,
		itemLeftMeta, itemIdentifier, itemColon, itemNumber, itemRightMeta, itemComma,
		itemLeftMeta, itemIdentifier, itemColon, itemNumber, itemRightMeta,
		itemRightBracket, itemRightMeta, itemEOF}

	JSON := "{myObjects:[{a:1},{b:2}]}"
	emittedTokens := lexJSON(JSON)
	validateTokens(expectedToken, emittedTokens, t)
}

func TestArrayOfArrays(t *testing.T) {
	expectedToken := []itemType{itemLeftMeta, itemIdentifier, itemColon,
		itemLeftBracket,
		itemLeftBracket, itemNumber, itemComma, itemNumber, itemRightBracket, itemComma,
		itemLeftBracket, itemNumber, itemComma, itemNumber, itemRightBracket,
		itemRightBracket,
		itemRightMeta, itemEOF}

	JSON := "{arrays:[[1,2],[3,4]]}"
	emittedTokens := lexJSON(JSON)
	validateTokens(expectedToken, emittedTokens, t)
}

func TestObjectWithinObject(t *testing.T) {
	expectedToken := []itemType{itemLeftMeta, itemIdentifier, itemColon, itemLeftMeta,
		itemIdentifier, itemColon, itemNumber, itemComma,
		itemIdentifier, itemColon, itemString,
		itemRightMeta, itemRightMeta, itemEOF}
	JSON := "{obj:{a:1,b:\"s\"}}"
	emittedTokens := lexJSON(JSON)
	validateTokens(expectedToken, emittedTokens, t)
}

func lexJSON(JSON string) []itemType {
	emittedTokens := make([]itemType, 0)

	_, tokens := lex("TestJSON", JSON)

	for t := range tokens {
		emittedTokens = append(emittedTokens, t.typ)
	}
	return emittedTokens
}
