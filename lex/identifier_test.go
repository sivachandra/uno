package lex

import (
	"testing"
	"uno/lex/token"
)

func TestIdentifierTokens(t *testing.T) {
	ts := token.NewTOISet([]token.Kind{
		token.KeywordClass,
		token.KeywordDef,
		token.Identifier,
	})

	err := matchTokens("test_data/identifiers_text", ts, []Token{
		Token{token.KeywordClass, "class", 1, 1},
		Token{token.Identifier, "MyClass", 1, 7},
		Token{token.KeywordDef, "def", 2, 1},
		Token{token.Identifier, "MyFunc", 2, 5},
		Token{token.KeywordDef, "def", 3, 1},
		Token{token.Identifier, "AnotherFunc", 3, 5},
		Token{token.Identifier, "_an_identifier", 5, 1},
		Token{token.Identifier, "_another_1_for_fun", 6, 1},
		Token{token.Identifier, "_take_100_then___", 6, 20},
	})

	if err != nil {
		t.Errorf(err.Error())
	}
}
