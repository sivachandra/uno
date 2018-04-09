package lex

import (
	"testing"
	"uno/lex/token"
)

// Individual token kinds like comments, identifiers, strings etc.
// are tested in corresponging test files. This file tests reading
// tokens from files having a combination of different kinds of
// tokens like in typical programming language file.

// This function does not really test the whole python
// language syntax but only a subset as a test for
// tokenizer.
func TestPythonText(t *testing.T) {
	ts := token.NewTOISet([]token.Kind{
		token.PySingleLineComment,
		token.KeywordClass,
		token.Identifier,
		token.LeftParen,
		token.RightParen,
		token.Colon,
		token.KeywordDef,
		token.Dot,
		token.Comma,
		token.ReturnArrow,
		token.KeywordReturn,
		token.LeftBracket,
		token.RightBracket,
		token.Assign,
		token.Indent,
		token.KeywordOr,
		token.SingleQuoteString,
	})

	err := matchTokens("test_data/python_text", ts, []Token{
		// Line 1
		Token{
			token.PySingleLineComment,
			"# This file has text data in the Python language syntax.",
			1,
			1,
		},
		// Line 2 is empty
		// Line 3
		Token{token.KeywordClass, "class", 3, 1},
		Token{token.Identifier, "MyClass", 3, 7},
		Token{token.LeftParen, "(", 3, 14},
		Token{token.Identifier, "object", 3, 15},
		Token{token.RightParen, ")", 3, 21},
		Token{token.Colon, ":", 3, 22},
		// Line 4
		Token{token.Indent, "    ", 4, 1},
		Token{token.KeywordDef, "def", 4, 5},
		Token{token.Identifier, "__init__", 4, 9},
		Token{token.LeftParen, "(", 4, 17},
		Token{token.Identifier, "self", 4, 18},
		Token{token.Comma, ",", 4, 22},
		Token{token.Identifier, "i", 4, 24},
		Token{token.Colon, ":", 4, 25},
		Token{token.Identifier, "int", 4, 27},
		Token{token.Comma, ",", 4, 30},
		Token{token.Identifier, "z", 4, 32},
		Token{token.Colon, ":", 4, 33},
		Token{token.Identifier, "str", 4, 35},
		Token{token.RightParen, ")", 4, 38},
		// Line 5
		Token{token.Indent, "        ", 5, 1},
		Token{token.Identifier, "self", 5, 9},
		Token{token.Dot, ".", 5, 13},
		Token{token.Identifier, "_i", 5, 14},
		Token{token.Assign, "=", 5, 17},
		Token{token.Identifier, "i", 5, 19},
		// Line 6
		Token{token.Indent, "        ", 6, 1},
		Token{token.Identifier, "self", 6, 9},
		Token{token.Dot, ".", 6, 13},
		Token{token.Identifier, "_z", 6, 14},
		Token{token.Assign, "=", 6, 17},
		Token{token.Identifier, "z", 6, 19},
		Token{token.KeywordOr, "or", 6, 21},
		Token{token.SingleQuoteString, "'\tHello, World'", 6, 24},
		// Line 7
		Token{token.Indent, "        ", 7, 1},
		Token{token.Identifier, "self", 7, 9},
		Token{token.Dot, ".", 7, 13},
		Token{token.Identifier, "_list", 7, 14},
		Token{token.Assign, "=", 7, 20},
		Token{token.LeftBracket, "[", 7, 22},
		Token{token.RightBracket, "]", 7, 23},
		Token{token.PySingleLineComment, "# Create an empty list.", 7, 26},
		// Line 8 is empty
		// Line 9
		Token{token.Indent, "    ", 9, 1},
		Token{token.KeywordDef, "def", 9, 5},
		Token{token.Identifier, "geti", 9, 9},
		Token{token.LeftParen, "(", 9, 13},
		Token{token.Identifier, "self", 9, 14},
		Token{token.RightParen, ")", 9, 18},
		Token{token.ReturnArrow, "->", 9, 20},
		Token{token.Identifier, "int", 9, 23},
		Token{token.Colon, ":", 9, 26},
		// Line 10
		Token{token.Indent, "        ", 10, 1},
		Token{token.KeywordReturn, "return", 10, 9},
		Token{token.Identifier, "self", 10, 16},
		Token{token.Dot, ".", 10, 20},
		Token{token.Identifier, "_i", 10, 21},
		// Line 11 is empty
		// Line 12
		Token{token.Indent, "    ", 12, 1},
		Token{token.KeywordDef, "def", 12, 5},
		Token{token.Identifier, "append", 12, 9},
		Token{token.LeftParen, "(", 12, 15},
		Token{token.Identifier, "self", 12, 16},
		Token{token.Comma, ",", 12, 20},
		Token{token.Identifier, "i", 12, 22},
		Token{token.Colon, ":", 12, 23},
		Token{token.Identifier, "int", 12, 25},
		Token{token.RightParen, ")", 12, 28},
		Token{token.Colon, ":", 12, 29},
		// Line 13
		Token{token.Indent, "        ", 13, 1},
		Token{token.KeywordReturn, "return", 13, 9},
		Token{token.Identifier, "self", 13, 16},
		Token{token.Dot, ".", 13, 20},
		Token{token.Identifier, "_list", 13, 21},
		Token{token.Dot, ".", 13, 26},
		Token{token.Identifier, "append", 13, 27},
		Token{token.LeftParen, "(", 13, 33},
		Token{token.Identifier, "i", 13, 34},
		Token{token.RightParen, ")", 13, 35},
		// Line 14 is empty
		// Line 15
		Token{token.Indent, "    ", 15, 1},
		Token{token.KeywordDef, "def", 15, 5},
		Token{token.Identifier, "get", 15, 9},
		Token{token.LeftParen, "(", 15, 12},
		Token{token.Identifier, "self", 15, 13},
		Token{token.Comma, ",", 15, 17},
		Token{token.Identifier, "index", 15, 19},
		Token{token.Colon, ":", 15, 24},
		Token{token.Identifier, "int", 15, 26},
		Token{token.RightParen, ")", 15, 29},
		Token{token.ReturnArrow, "->", 15, 31},
		Token{token.Identifier, "int", 15, 34},
		Token{token.Colon, ":", 15, 37},
		// Line 16
		Token{token.Indent, "        ", 16, 1},
		Token{token.KeywordReturn, "return", 16, 9},
		Token{token.Identifier, "self", 16, 16},
		Token{token.Dot, ".", 16, 20},
		Token{token.Identifier, "_list", 16, 21},
		Token{token.LeftBracket, "[", 16, 26},
		Token{token.Identifier, "index", 16, 27},
		Token{token.RightBracket, "]", 16, 32},
	})

	if err != nil {
		t.Errorf(err.Error())
	}
}
