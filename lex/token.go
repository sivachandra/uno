package lex

type Token struct {
	Kind  uint32
	Value string
	Line  uint32
	Col   uint32
}

func newToken(tt uint32, val []rune, l uint32, c uint32) *Token {
	t := new(Token)
	t.Kind = tt
	t.Line = l
	t.Col = c
	t.Value = string(val)
	return t
}

