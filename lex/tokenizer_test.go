package lex

import (
	"testing"
	"uno/lex/token_kind"
)

// Individual token kinds like comments, identifiers, strings etc.
// are tested in corresponging test files. This file tests reading
// tokens from files having a combination of different kinds of
// tokens like in typical programming language file.

// This function does not really test the whole python
// language syntax but only a subset as a test for
// tokenizer.
func TestPythonText(t *testing.T) {
	ts := NewTokenKindSet([]uint32{
		token_kind.PySingleLineComment,
		token_kind.KeywordClass,
		token_kind.Identifier,
		token_kind.LeftParen,
		token_kind.RightParen,
		token_kind.Colon,
		token_kind.KeywordDef,
		token_kind.Dot,
		token_kind.Comma,
		token_kind.ReturnArrow,
		token_kind.KeywordReturn,
		token_kind.LeftBracket,
		token_kind.RightBracket,
		token_kind.Assign,
		token_kind.Indent,
		token_kind.KeywordOr,
		token_kind.SingleQuoteString,
	})

	err := matchTokens("test_data/python_text", ts, []Token{
		// Line 1
		Token{
			token_kind.PySingleLineComment,
			"# This file has text data in the Python language syntax.",
			1,
			1,
		},
		// Line 2 is empty
		// Line 3
		Token{token_kind.KeywordClass, "class", 3, 1},
		Token{token_kind.Identifier, "MyClass", 3, 7},
		Token{token_kind.LeftParen, "(", 3, 14},
		Token{token_kind.Identifier, "object", 3, 15},
		Token{token_kind.RightParen, ")", 3, 21},
		Token{token_kind.Colon, ":", 3, 22},
		// Line 4
		Token{token_kind.Indent, "    ", 4, 1},
		Token{token_kind.KeywordDef, "def", 4, 5},
		Token{token_kind.Identifier, "__init__", 4, 9},
		Token{token_kind.LeftParen, "(", 4, 17},
		Token{token_kind.Identifier, "self", 4, 18},
		Token{token_kind.Comma, ",", 4, 22},
		Token{token_kind.Identifier, "i", 4, 24},
		Token{token_kind.Colon, ":", 4, 25},
		Token{token_kind.Identifier, "int", 4, 27},
		Token{token_kind.Comma, ",", 4, 30},
		Token{token_kind.Identifier, "z", 4, 32},
		Token{token_kind.Colon, ":", 4, 33},
		Token{token_kind.Identifier, "str", 4, 35},
		Token{token_kind.RightParen, ")", 4, 38},
		// Line 5
		Token{token_kind.Indent, "        ", 5, 1},
		Token{token_kind.Identifier, "self", 5, 9},
		Token{token_kind.Dot, ".", 5, 13},
		Token{token_kind.Identifier, "_i", 5, 14},
		Token{token_kind.Assign, "=", 5, 17},
		Token{token_kind.Identifier, "i", 5, 19},
		// Line 6
		Token{token_kind.Indent, "        ", 6, 1},
		Token{token_kind.Identifier, "self", 6, 9},
		Token{token_kind.Dot, ".", 6, 13},
		Token{token_kind.Identifier, "_z", 6, 14},
		Token{token_kind.Assign, "=", 6, 17},
		Token{token_kind.Identifier, "z", 6, 19},
		Token{token_kind.KeywordOr, "or", 6, 21},
		Token{token_kind.SingleQuoteString, "'\tHello, World'", 6, 24},
		// Line 7
		Token{token_kind.Indent, "        ", 7, 1},
		Token{token_kind.Identifier, "self", 7, 9},
		Token{token_kind.Dot, ".", 7, 13},
		Token{token_kind.Identifier, "_list", 7, 14},
		Token{token_kind.Assign, "=", 7, 20},
		Token{token_kind.LeftBracket, "[", 7, 22},
		Token{token_kind.RightBracket, "]", 7, 23},
		Token{token_kind.PySingleLineComment, "# Create an empty list.", 7, 26},
		// Line 8 is empty
		// Line 9
		Token{token_kind.Indent, "    ", 9, 1},
		Token{token_kind.KeywordDef, "def", 9, 5},
		Token{token_kind.Identifier, "geti", 9, 9},
		Token{token_kind.LeftParen, "(", 9, 13},
		Token{token_kind.Identifier, "self", 9, 14},
		Token{token_kind.RightParen, ")", 9, 18},
		Token{token_kind.ReturnArrow, "->", 9, 20},
		Token{token_kind.Identifier, "int", 9, 23},
		Token{token_kind.Colon, ":", 9, 26},
		// Line 10
		Token{token_kind.Indent, "        ", 10, 1},
		Token{token_kind.KeywordReturn, "return", 10, 9},
		Token{token_kind.Identifier, "self", 10, 16},
		Token{token_kind.Dot, ".", 10, 20},
		Token{token_kind.Identifier, "_i", 10, 21},
		// Line 11 is empty
		// Line 12
		Token{token_kind.Indent, "    ", 12, 1},
		Token{token_kind.KeywordDef, "def", 12, 5},
		Token{token_kind.Identifier, "append", 12, 9},
		Token{token_kind.LeftParen, "(", 12, 15},
		Token{token_kind.Identifier, "self", 12, 16},
		Token{token_kind.Comma, ",", 12, 20},
		Token{token_kind.Identifier, "i", 12, 22},
		Token{token_kind.Colon, ":", 12, 23},
		Token{token_kind.Identifier, "int", 12, 25},
		Token{token_kind.RightParen, ")", 12, 28},
		Token{token_kind.Colon, ":", 12, 29},
		// Line 13
		Token{token_kind.Indent, "        ", 13, 1},
		Token{token_kind.KeywordReturn, "return", 13, 9},
		Token{token_kind.Identifier, "self", 13, 16},
		Token{token_kind.Dot, ".", 13, 20},
		Token{token_kind.Identifier, "_list", 13, 21},
		Token{token_kind.Dot, ".", 13, 26},
		Token{token_kind.Identifier, "append", 13, 27},
		Token{token_kind.LeftParen, "(", 13, 33},
		Token{token_kind.Identifier, "i", 13, 34},
		Token{token_kind.RightParen, ")", 13, 35},
		// Line 14 is empty
		// Line 15
		Token{token_kind.Indent, "    ", 15, 1},
		Token{token_kind.KeywordDef, "def", 15, 5},
		Token{token_kind.Identifier, "get", 15, 9},
		Token{token_kind.LeftParen, "(", 15, 12},
		Token{token_kind.Identifier, "self", 15, 13},
		Token{token_kind.Comma, ",", 15, 17},
		Token{token_kind.Identifier, "index", 15, 19},
		Token{token_kind.Colon, ":", 15, 24},
		Token{token_kind.Identifier, "int", 15, 26},
		Token{token_kind.RightParen, ")", 15, 29},
		Token{token_kind.ReturnArrow, "->", 15, 31},
		Token{token_kind.Identifier, "int", 15, 34},
		Token{token_kind.Colon, ":", 15, 37},
		// Line 16
		Token{token_kind.Indent, "        ", 16, 1},
		Token{token_kind.KeywordReturn, "return", 16, 9},
		Token{token_kind.Identifier, "self", 16, 16},
		Token{token_kind.Dot, ".", 16, 20},
		Token{token_kind.Identifier, "_list", 16, 21},
		Token{token_kind.LeftBracket, "[", 16, 26},
		Token{token_kind.Identifier, "index", 16, 27},
		Token{token_kind.RightBracket, "]", 16, 32},
	})

	if err != nil {
		t.Errorf(err.Error())
	}
}
