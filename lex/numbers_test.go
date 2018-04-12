package lex

import (
	"testing"
	"uno/lex/token_kind"
)

func TestNumberTokens(t *testing.T) {
	ts := NewTokenKindSet([]uint32{
		token_kind.DecimalInteger,
		token_kind.HexInteger,
		token_kind.OctInteger,
		token_kind.FloatNumber,
	})

	err := matchTokens("test_data/numbers_text", ts, []Token{
		Token{token_kind.DecimalInteger, "12345", 1, 1},
		Token{token_kind.OctInteger, "012345", 1, 7},
		Token{token_kind.HexInteger, "0x12345", 1, 14},
		Token{token_kind.HexInteger, "0xabcdef", 2, 1},
		Token{token_kind.HexInteger, "0xbadface", 3, 1},
		Token{token_kind.DecimalInteger, "0", 4, 1},
		Token{token_kind.DecimalInteger, "0", 4, 3},
		Token{token_kind.DecimalInteger, "0", 4, 5},
		Token{token_kind.OctInteger, "0123", 5, 1},
		Token{token_kind.FloatNumber, "123.456", 5, 6},
		Token{token_kind.FloatNumber, "1.2e15", 5, 14},
		Token{token_kind.DecimalInteger, "54321", 5, 21},
		Token{token_kind.FloatNumber, "5678e6", 5, 27},
		Token{token_kind.FloatNumber, ".432", 6, 1},
		Token{token_kind.FloatNumber, "0.234", 6, 6},
		Token{token_kind.FloatNumber, "0.e12", 6, 12},
		Token{token_kind.FloatNumber, "0.1e13", 6, 18},
	})

	if err != nil {
		t.Errorf(err.Error())
	}
}
