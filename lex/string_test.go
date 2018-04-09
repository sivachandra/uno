package lex

import (
	"testing"
	"uno/lex/token"
)

func TestWithSingleQuoteCharacter(t *testing.T) {
	ts := token.NewTOISet([]token.Kind{
		// We include double quote strings as they can
		// exist along with single quote characters.
		token.DoubleQuoteString,
		token.SingleQuoteCharacter,
	})

	err := matchTokens("test_data/single_quote_char_text", ts, []Token{
		Token{token.SingleQuoteCharacter, "'a'", 1, 1},
		Token{token.SingleQuoteCharacter, "'\n'", 1, 5},
		Token{token.SingleQuoteCharacter, "'\t'", 1, 10},
		Token{token.DoubleQuoteString, "\"hello, world\"", 2, 1},
		Token{token.DoubleQuoteString, "\"hello,\nworld\"", 3, 1},
		Token{token.DoubleQuoteString, "\"hello,\tworld\"", 4, 1},
		Token{token.DoubleQuoteString, "\"hello,\vworld\"", 4, 17},
	})

	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestWithSingleQuoteString(t *testing.T) {
	ts := token.NewTOISet([]token.Kind{
		// We include double quote strings as they can
		// exist along with single quote strings.
		token.DoubleQuoteString,
		token.SingleQuoteString,
	})

	err := matchTokens("test_data/single_quote_string_text", ts, []Token{
		Token{token.SingleQuoteString, "'hello, world'", 1, 1},
		Token{token.SingleQuoteString, "'hello,\nworld'", 2, 1},
		Token{token.SingleQuoteString, "'hello,\tworld'", 3, 1},
		Token{token.SingleQuoteString, "'hello,\vworld'", 3, 17},
		Token{token.DoubleQuoteString, "\"hello, world\"", 5, 1},
		Token{token.DoubleQuoteString, "\"hello,\nworld\"", 6, 1},
		Token{token.DoubleQuoteString, "\"hello,\tworld\"", 7, 1},
		Token{token.DoubleQuoteString, "\"hello,\vworld\"", 7, 17},
	})

	if err != nil {
		t.Errorf(err.Error())
	}
}
