package lex

import (
	"testing"
	"uno/lex/token"
)

func TestNumberTokens(t *testing.T) {
	ts := token.NewTOISet([]token.Kind{
		token.DecimalInteger,
		token.HexInteger,
		token.OctInteger,
		token.FloatNumber,
	})

	err := matchTokens("test_data/numbers_text", ts, []Token{
		Token{token.DecimalInteger, "12345", 1, 1},
		Token{token.OctInteger, "012345", 1, 7},
		Token{token.HexInteger, "0x12345", 1, 14},
		Token{token.HexInteger, "0xabcdef", 2, 1},
		Token{token.HexInteger, "0xbadface", 3, 1},
		Token{token.DecimalInteger, "0", 4, 1},
		Token{token.DecimalInteger, "0", 4, 3},
		Token{token.DecimalInteger, "0", 4, 5},
		Token{token.OctInteger, "0123", 5, 1},
		Token{token.FloatNumber, "123.456", 5, 6},
		Token{token.FloatNumber, "1.2e15", 5, 14},
		Token{token.DecimalInteger, "54321", 5, 21},
		Token{token.FloatNumber, "5678e6", 5, 27},
		Token{token.FloatNumber, ".432", 6, 1},
		Token{token.FloatNumber, "0.234", 6, 6},
		Token{token.FloatNumber, "0.e12", 6, 12},
		Token{token.FloatNumber, "0.1e13", 6, 18},
	})

	if err != nil {
		t.Errorf(err.Error())
	}
}
