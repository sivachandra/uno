package lex

import (
	"testing"
	"uno/lex/token_kind"
)

func TestWithSingleQuoteCharacter(t *testing.T) {
	ts := NewTokenKindSet([]uint32{
		// We include double quote strings as they can
		// exist along with single quote characters.
		token_kind.DoubleQuoteString,
		token_kind.SingleQuoteCharacter,
	})

	err := matchTokens("test_data/single_quote_char_text", ts, []Token{
		Token{token_kind.SingleQuoteCharacter, "'a'", 1, 1},
		Token{token_kind.SingleQuoteCharacter, "'\n'", 1, 5},
		Token{token_kind.SingleQuoteCharacter, "'\t'", 1, 10},
		Token{token_kind.DoubleQuoteString, "\"hello, world\"", 2, 1},
		Token{token_kind.DoubleQuoteString, "\"hello,\nworld\"", 3, 1},
		Token{token_kind.DoubleQuoteString, "\"hello,\tworld\"", 4, 1},
		Token{token_kind.DoubleQuoteString, "\"hello,\vworld\"", 4, 17},
	})

	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestWithSingleQuoteString(t *testing.T) {
	ts := NewTokenKindSet([]uint32{
		// We include double quote strings as they can
		// exist along with single quote strings.
		token_kind.DoubleQuoteString,
		token_kind.SingleQuoteString,
	})

	err := matchTokens("test_data/single_quote_string_text", ts, []Token{
		Token{token_kind.SingleQuoteString, "'hello, world'", 1, 1},
		Token{token_kind.SingleQuoteString, "'hello,\nworld'", 2, 1},
		Token{token_kind.SingleQuoteString, "'hello,\tworld'", 3, 1},
		Token{token_kind.SingleQuoteString, "'hello,\vworld'", 3, 17},
		Token{token_kind.DoubleQuoteString, "\"hello, world\"", 5, 1},
		Token{token_kind.DoubleQuoteString, "\"hello,\nworld\"", 6, 1},
		Token{token_kind.DoubleQuoteString, "\"hello,\tworld\"", 7, 1},
		Token{token_kind.DoubleQuoteString, "\"hello,\vworld\"", 7, 17},
	})

	if err != nil {
		t.Errorf(err.Error())
	}
}
