package lex

import (
	"fmt"
	"io"
	"uno/lex/char"
	"uno/lex/token"
)

func (tz *Tokenizer) readPythonStyleComment() (*Token, error) {
	line := tz.r.NextLine()
	col := tz.r.NextCol()

	h, err := tz.r.ReadChar()
	if err != nil {
		return nil, err
	}
	if err != nil {
		err = fmt.Errorf("Expected to read a comment begging with '#'.\n%s", err.Error())
		return nil, err
	}

	var s []rune
	c := h
	for c != char.NewLine {
		s = append(s, c)
		c, err = tz.r.ReadChar()
		if err == io.EOF {
			break
		}
	}

	return newToken(token.PySingleLineComment, s, line, col), nil
}

func (tz *Tokenizer) readCStyleSingleLineComment() (*Token, error) {
	line := tz.r.NextLine()
	col := tz.r.NextCol()

	ss, err := tz.r.ReadSlice(2)
	if err != nil {
		err = fmt.Errorf("Expected to read a comment begging with '//'.\n%s", err.Error())
		return nil, err
	}

	if string(ss) != "//" {
		return nil, fmt.Errorf("Expected a comment beginning with '//'.")
	}

	var s []rune
	s = append(s, ss[0])
	c := ss[1]
	for c != char.NewLine {
		s = append(s, c)
		c, err = tz.r.ReadChar()
		if err == io.EOF {
			break
		}
	}

	return newToken(token.CSingleLineComment, s, line, col), nil
}

func (tz *Tokenizer) readCStyleMultiLineComment() (*Token, error) {
	line := tz.r.NextLine()
	col := tz.r.NextCol()

	ss, err := tz.r.ReadSlice(2)
	if err != nil {
		err = fmt.Errorf("Expected to read a comment begging with '/*'.\n%s", err.Error())
		return nil, err
	}

	if string(ss) != "/*" {
		return nil, fmt.Errorf("Expected a comment beginning with '/*'.")
	}

	var s []rune
	s = append(s, ss...)
	for true {
		c, err := tz.r.ReadChar()
		if err != nil {
			return nil, fmt.Errorf("Error reading multine comment.\n%s", err.Error())
		}

		s = append(s, c)

		if l := len(s); string(s[l-2:l]) == "*/" {
			break
		}
	}

	return newToken(token.CMultiLineComment, s, line, col), nil
}
