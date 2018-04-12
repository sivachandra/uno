package lex

import (
	"fmt"
	"uno/lex/token_kind"
)

type GoESR struct {
}

var goCommonEscSeq = map[rune]rune{
	'a':  '\a',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'v':  '\v',
	'\\': '\\',
}

func (esr GoESR) ReadChar(r *CharReader, tt uint32) (rune, error) {
	c, err := r.ReadChar()
	if err != nil {
		return rune(0), err
	}
	val, valid := goCommonEscSeq[c]
	if valid {
		return val, nil
	}

	switch tt {
	case token_kind.SingleQuoteCharacter:
		if c == '\'' {
			return c, nil
		}
	case token_kind.DoubleQuoteString:
		if c == '"' {
			return c, nil
		}
	}
	return 0, fmt.Errorf("Invalid escape sequence.")
}
