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
	expectedToken := []itemType{itemLeftMeta, itemRightMeta}
	emittedTokens := lexJSON("{}")
	validateTokens(expectedToken, emittedTokens, t)
}

func TestLexEmptyObjects(t *testing.T) {
	expectedToken := []itemType{itemLeftMeta,
		itemIdentifier, itemLeftMeta,
		itemIdentifier, itemLeftMeta,
		itemIdentifier, itemLeftMeta,
		itemIdentifier, itemLeftMeta,
		itemRightMeta, itemRightMeta, itemRightMeta, itemRightMeta, itemRightMeta}
	emittedTokens := lexJSON("{a:{b:{c:{d:{}}}}}")
	validateTokens(expectedToken, emittedTokens, t)
}

func TestLexEmptyArray(t *testing.T) {
	expectedToken := []itemType{itemLeftMeta, itemIdentifier, itemLeftBracket, itemRightBracket, itemRightMeta}
	emittedTokens := lexJSON("{a:[]}")
	validateTokens(expectedToken, emittedTokens, t)
}

func TestLexEmptyArrays(t *testing.T) {
	expectedToken := []itemType{itemLeftMeta, itemIdentifier,
		itemLeftBracket, itemLeftBracket, itemLeftBracket, itemLeftBracket, itemLeftBracket,
		itemRightBracket, itemRightBracket, itemRightBracket, itemRightBracket, itemRightBracket,
		itemRightMeta}
	emittedTokens := lexJSON("{a:[[[[[]]]]]}")
	validateTokens(expectedToken, emittedTokens, t)
}

func TestLexNumericIdentifier(t *testing.T) {
	expectedToken := []itemType{itemLeftMeta, itemIdentifier, itemNumber, itemRightMeta}
	JSON := "{a:1}"
	emittedTokens := lexJSON(JSON)
	validateTokens(expectedToken, emittedTokens, t)
}

func TestLexStringIdentifier(t *testing.T) {
	expectedToken := []itemType{itemLeftMeta, itemIdentifier, itemString, itemRightMeta}
	JSON := "{b:\"value\"}"
	emittedTokens := lexJSON(JSON)
	validateTokens(expectedToken, emittedTokens, t)
}

func TestLexTwoIdentifiers(t *testing.T) {
	expectedToken := []itemType{itemLeftMeta, itemIdentifier, itemNumber, itemComma, itemIdentifier, itemString, itemRightMeta}
	JSON := "{a:1,b:\"value\"}"
	emittedTokens := lexJSON(JSON)
	validateTokens(expectedToken, emittedTokens, t)
}

func TestLexArrayIdentifier(t *testing.T) {
	expectedToken := []itemType{itemLeftMeta, itemIdentifier, itemLeftBracket, itemNumber, itemComma, itemNumber, itemComma, itemNumber, itemRightBracket, itemRightMeta}
	JSON := "{a:[1,2,3]}"
	emittedTokens := lexJSON(JSON)
	validateTokens(expectedToken, emittedTokens, t)
}

func TestLexArrayOfObjects(t *testing.T) {
	expectedToken := []itemType{itemLeftMeta, itemIdentifier, itemLeftBracket,
		itemLeftMeta, itemIdentifier, itemNumber, itemRightMeta, itemComma,
		itemLeftMeta, itemIdentifier, itemNumber, itemRightMeta,
		itemRightBracket, itemRightMeta}

	JSON := "{myObjects:[{a:1},{b:2}]}"
	emittedTokens := lexJSON(JSON)
	validateTokens(expectedToken, emittedTokens, t)
}

func TestArrayOfArrays(t *testing.T) {
	expectedToken := []itemType{itemLeftMeta, itemIdentifier,
		itemLeftBracket,
		itemLeftBracket, itemNumber, itemComma, itemNumber, itemRightBracket, itemComma,
		itemLeftBracket, itemNumber, itemComma, itemNumber, itemRightBracket,
		itemRightBracket,
		itemRightMeta}

	JSON := "{arrays:[[1,2],[3,4]]}"
	emittedTokens := lexJSON(JSON)
	validateTokens(expectedToken, emittedTokens, t)
}

func TestObjectWithinObject(t *testing.T) {
	expectedToken := []itemType{itemLeftMeta, itemIdentifier, itemLeftMeta,
		itemIdentifier, itemNumber, itemComma,
		itemIdentifier, itemString,
		itemRightMeta, itemRightMeta}
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
