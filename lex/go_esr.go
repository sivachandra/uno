package lex

import (
	"fmt"
	"uno/lex/token"
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

func (esr GoESR) ReadChar(r *CharReader, tt token.Kind) (rune, error) {
	c, err := r.ReadChar()
	if err != nil {
		return rune(0), err
	}
	val, valid := goCommonEscSeq[c]
	if valid {
		return val, nil
	}

	switch tt {
	case token.SingleQuoteCharacter:
		if c == '\'' {
			return c, nil
		}
	case token.DoubleQuoteString:
		if c == '"' {
			return c, nil
		}
	}
	return 0, fmt.Errorf("Invalid escape sequence.")
}
