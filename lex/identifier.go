package lex

import (
	"fmt"
	"unicode"
	"uno/lex/char"
	"uno/lex/token_kind"
)

var KeywordMap = map[string]uint32{
	"and":      token_kind.KeywordAnd,
	"as":       token_kind.KeywordAs,
	"assert":   token_kind.KeywordAssert,
	"break":    token_kind.KeywordBreak,
	"class":    token_kind.KeywordClass,
	"const":    token_kind.KeywordConst,
	"continue": token_kind.KeywordContinue,
	"def":      token_kind.KeywordDef,
	"del":      token_kind.KeywordDel,
	"elif":     token_kind.KeywordElif,
	"else":     token_kind.KeywordElse,
	"except":   token_kind.KeywordExcept,
	"false":    token_kind.KeywordCFalse,
	"False":    token_kind.KeywordPyFalse,
	"finally":  token_kind.KeywordFinally,
	"for":      token_kind.KeywordFor,
	"from":     token_kind.KeywordFrom,
	"global":   token_kind.KeywordGlobal,
	"if":       token_kind.KeywordIf,
	"import":   token_kind.KeywordImport,
	"in":       token_kind.KeywordIn,
	"is":       token_kind.KeywordIs,
	"lambda":   token_kind.KeywordLambda,
	"not":      token_kind.KeywordNot,
	"null":     token_kind.KeywordNull,
	"or":       token_kind.KeywordOr,
	"pass":     token_kind.KeywordPass,
	"raise":    token_kind.KeywordRaise,
	"return":   token_kind.KeywordReturn,
	"true":     token_kind.KeywordCTrue,
	"True":     token_kind.KeywordPyTrue,
	"try":      token_kind.KeywordTry,
	"while":    token_kind.KeywordWhile,
	"with":     token_kind.KeywordWith,
	"yield":    token_kind.KeywordYield,
}

func isIdentifierBeginChar(c rune) bool {
	return c == char.Underscore || unicode.IsLetter(rune(c))
}

func isIdentifierContinuationChar(c rune) bool {
	return isIdentifierBeginChar(c) || unicode.IsNumber(rune(c))
}

func (tz *Tokenizer) readIdentifierString() ([]rune, error) {
	c, err := tz.r.PeekChar()
	if err != nil {
		return nil, fmt.Errorf("Error reading identifier.\n%s", err.Error())
	}
	if !isIdentifierBeginChar(c) {
		return nil, fmt.Errorf("Invalid identifier begin character '%c'.", c)
	}

	c, err = tz.r.ReadChar()
	if err != nil {
		return nil, fmt.Errorf("Error reading identifier.\n%s", err.Error())
	}

	var id []rune
	id = append(id, c)

	for true {
		c, err = tz.r.PeekChar()
		if err != nil {
			return nil, fmt.Errorf("Error reading identifier.\n%s", err.Error())
		}
		if !isIdentifierContinuationChar(c) {
			break
		}

		// Actually read out the character.
		c, err = tz.r.ReadChar()
		if err != nil {
			return nil, fmt.Errorf("Error reading identifier.\n%s", err.Error())
		}

		id = append(id, c)
	}

	return id, nil
}

func (tz *Tokenizer) readIdentifier() (*Token, error) {
	col := tz.r.NextCol()

	id, err := tz.readIdentifierString()
	if err != nil {
		return nil, err
	}

	tt, e := KeywordMap[string(id)]
	if e && tz.ts.Contains(tt) {
		return newToken(tt, id, tz.r.Line(), col), nil
	} else {
		return newToken(token_kind.Identifier, id, tz.r.Line(), col), nil
	}
}
