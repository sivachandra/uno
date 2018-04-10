package lex

import (
	"uno/lex/token"
)

type Token struct {
	Kind  token.Kind
	Value string
	Line  uint32
	Col   uint32
}

func newToken(tt token.Kind, val []rune, l uint32, c uint32) *Token {
	t := new(Token)
	t.Kind = tt
	t.Line = l
	t.Col = c
	t.Value = string(val)
	return t
}

