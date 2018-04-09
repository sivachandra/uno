package lex

import (
	"testing"
	"uno/lex/token"
)

func TestCommentTokens(t *testing.T) {
	// Though this is a test for comments, we include
	// token.Identifier so that we can test comments at the
	// end of line in this form:
	//   an_indentifier  # My comment end of the line comment
	ts := token.NewTOISet([]token.Kind{
		token.PySingleLineComment,
		token.CSingleLineComment,
		token.CMultiLineComment,
		token.Identifier,
	})

	mlComment1 := "/* C style multiline comment\n                extending to a new line.  */"
	mlComment2 := "/* Another C style multiline comment\n   extending to a new line.  */"

	err := matchTokens("test_data/comments_text", ts, []Token{
		Token{
			token.PySingleLineComment,
			"# This file has simple identifier tokens and comments.",
			1,
			1,
		},
		Token{
			token.PySingleLineComment,
			"# It is used for testing of reading comment tokens.",
			2,
			1,
		},
		Token{token.Identifier, "identifier1", 4, 1},
		Token{token.PySingleLineComment, "# A single line python comment", 4, 14},
		Token{token.Identifier, "identifier2", 5, 1},
		Token{token.CSingleLineComment, "// A C++ style single line comment.", 5, 14},
		Token{token.Identifier, "identifier3", 6, 1},
		Token{token.CMultiLineComment, mlComment1, 6, 14},
		Token{token.CMultiLineComment, mlComment2, 9, 1},
		Token{token.Identifier, "identifier3", 11, 1},
	})

	if err != nil {
		t.Errorf(err.Error())
	}
}
