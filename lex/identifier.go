package lex

import (
	"fmt"
	"unicode"
	"uno/lex/char"
	"uno/lex/token"
)

var Keywords = map[string]token.Kind{
	"and":      token.KeywordAnd,
	"as":       token.KeywordAs,
	"assert":   token.KeywordAssert,
	"break":    token.KeywordBreak,
	"class":    token.KeywordClass,
	"const":    token.KeywordConst,
	"continue": token.KeywordContinue,
	"def":      token.KeywordDef,
	"del":      token.KeywordDel,
	"elif":     token.KeywordElif,
	"else":     token.KeywordElse,
	"except":   token.KeywordExcept,
	"false":    token.KeywordCFalse,
	"False":    token.KeywordPyFalse,
	"finally":  token.KeywordFinally,
	"for":      token.KeywordFor,
	"from":     token.KeywordFrom,
	"global":   token.KeywordGlobal,
	"if":       token.KeywordIf,
	"import":   token.KeywordImport,
	"in":       token.KeywordIn,
	"is":       token.KeywordIs,
	"lambda":   token.KeywordLambda,
	"not":      token.KeywordNot,
	"null":     token.KeywordNull,
	"or":       token.KeywordOr,
	"pass":     token.KeywordPass,
	"raise":    token.KeywordRaise,
	"return":   token.KeywordReturn,
	"true":     token.KeywordCTrue,
	"True":     token.KeywordPyTrue,
	"try":      token.KeywordTry,
	"while":    token.KeywordWhile,
	"with":     token.KeywordWith,
	"yield":    token.KeywordYield,
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

	tt, e := Keywords[string(id)]
	if e && tz.ts.Contains(tt) {
		return newToken(tt, id, tz.r.Line(), col), nil
	} else {
		return newToken(token.Identifier, id, tz.r.Line(), col), nil
	}
}
