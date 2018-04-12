package lex

import (
	"testing"
	"uno/lex/token_kind"
)

func TestIdentifierTokens(t *testing.T) {
	ts := NewTokenKindSet([]uint32{
		token_kind.KeywordClass,
		token_kind.KeywordDef,
		token_kind.Identifier,
	})

	err := matchTokens("test_data/identifiers_text", ts, []Token{
		Token{token_kind.KeywordClass, "class", 1, 1},
		Token{token_kind.Identifier, "MyClass", 1, 7},
		Token{token_kind.KeywordDef, "def", 2, 1},
		Token{token_kind.Identifier, "MyFunc", 2, 5},
		Token{token_kind.KeywordDef, "def", 3, 1},
		Token{token_kind.Identifier, "AnotherFunc", 3, 5},
		Token{token_kind.Identifier, "_an_identifier", 5, 1},
		Token{token_kind.Identifier, "_another_1_for_fun", 6, 1},
		Token{token_kind.Identifier, "_take_100_then___", 6, 20},
	})

	if err != nil {
		t.Errorf(err.Error())
	}
}
