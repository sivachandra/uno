package lex

import (
	"testing"
	"uno/lex/token_kind"
)

func TestCommentTokens(t *testing.T) {
	// Though this is a test for comments, we include
	// token_kind.Identifier so that we can test comments at the
	// end of line in this form:
	//   an_indentifier  # My comment end of the line comment
	ts := NewTokenKindSet([]uint32{
		token_kind.PySingleLineComment,
		token_kind.CSingleLineComment,
		token_kind.CMultiLineComment,
		token_kind.Identifier,
	})

	mlComment1 := "/* C style multiline comment\n                extending to a new line.  */"
	mlComment2 := "/* Another C style multiline comment\n   extending to a new line.  */"

	err := matchTokens("test_data/comments_text", ts, []Token{
		Token{
			token_kind.PySingleLineComment,
			"# This file has simple identifier tokens and comments.",
			1,
			1,
		},
		Token{
			token_kind.PySingleLineComment,
			"# It is used for testing of reading comment tokens.",
			2,
			1,
		},
		Token{token_kind.Identifier, "identifier1", 4, 1},
		Token{token_kind.PySingleLineComment, "# A single line python comment", 4, 14},
		Token{token_kind.Identifier, "identifier2", 5, 1},
		Token{token_kind.CSingleLineComment, "// A C++ style single line comment.", 5, 14},
		Token{token_kind.Identifier, "identifier3", 6, 1},
		Token{token_kind.CMultiLineComment, mlComment1, 6, 14},
		Token{token_kind.CMultiLineComment, mlComment2, 9, 1},
		Token{token_kind.Identifier, "identifier3", 11, 1},
	})

	if err != nil {
		t.Errorf(err.Error())
	}
}
