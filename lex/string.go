package lex

import (
	"fmt"
	"uno/lex/char"
	"uno/lex/token"
)

type EscSeqReader interface {
	// Read the escape character after the '\' character.
	//
	// tt is the token type in which this escape sequence occurs.
	//
	// Returns the read character.
	ReadChar(r *CharReader, tt token.Kind) (rune, error)
}

func isSpace(c rune) bool {
	return c == char.Space || c == char.Tab
}

func isAnyWhiteSpace(c rune) bool {
	return isSpace(c) || c == char.NewLine || c == char.Return
}

func (tz *Tokenizer) readIndentToken() (*Token, error) {
	line := tz.r.NextLine()
	col := tz.r.NextCol()

	var v []rune
	for true {
		c, e := tz.r.PeekChar()
		if e != nil {
			return nil, fmt.Errorf("Error reading indent.\n%s", e.Error())
		}
		if isSpace(c) {
			c, e := tz.r.ReadChar()
			if e != nil {
				return nil, fmt.Errorf("Error reading indent.\n%s", e.Error())
			}

			v = append(v, c)
		} else {
			break
		}
	}

	return newToken(token.Indent, v, line, col), nil
}

func (tz *Tokenizer) skipSpace() error {
	for true {
		c, e := tz.r.PeekChar()
		if e != nil {
			return fmt.Errorf("Error skipping space.\n%s", e.Error())
		}
		if isSpace(c) {
			_, e := tz.r.ReadChar()
			if e != nil {
				return fmt.Errorf("Error skipping space.\n%s", e.Error())
			}
		} else {
			break
		}
	}
	return nil
}

func (tz *Tokenizer) readSingleQuoteCharacter() (*Token, error) {
	// Save the starting column.
	line := tz.r.NextLine()
	col := tz.r.NextCol()

	// A sigle quote character is of the form '<c>' or '\<c>'.
	q, err := tz.r.ReadChar()
	if err != nil {
		return nil, fmt.Errorf("Error reading single quote character.\n%s", err.Error())
	}
	if q != char.SingleQuote {
		return nil, fmt.Errorf("Trying to read single quote, but found '%s'.", string(q))
	}

	c, err := tz.r.ReadChar()
	if err != nil {
		return nil, fmt.Errorf("Error reading single quote character.\n%s", err.Error())
	}
	// If the second char is not an escape char, then it should
	// not be a newline char and the third char should be the closing
	// single quote.
	if c != '\\' {
		if c == char.NewLine || c == char.Return {
			return nil, fmt.Errorf("Invalid newline after single quote.")
		}

		q, err = tz.r.ReadChar()
		if q != char.SingleQuote {
			return nil, fmt.Errorf("Missing/incorrectly placed closing single quote.")
		}

		t := newToken(
			token.SingleQuoteCharacter,
			[]rune{char.SingleQuote, c, char.SingleQuote},
			line, col)
		return t, nil
	}

	// If the second char is a '\', then it is an escape sequence to be read
	// according to the language specific rules.
	c, err = tz.esr.ReadChar(tz.r, token.SingleQuoteCharacter)
	if err != nil {
		return nil, err
	}

	q, err = tz.r.ReadChar()
	if err != nil {
		return nil, fmt.Errorf(
			"Error reading terminating single quote character.\n%s",
			err.Error())
	}
	if q != char.SingleQuote {
		return nil, fmt.Errorf("Missing or incorrectly placed closing single quote.")
	}

	t := newToken(
		token.SingleQuoteCharacter,
		[]rune{char.SingleQuote, c, char.SingleQuote},
		line, col)
	return t, nil
}

var quoteTokenKind = map[rune]token.Kind {
	char.DoubleQuote: token.DoubleQuoteString,
	char.SingleQuote: token.SingleQuoteString,
	char.BackQuote:   token.BackQuoteString,
}

func (tz *Tokenizer) readQuotedString(raw bool) (*Token, error) {
	// Save the beginning line and column for reporting.
	col := tz.r.NextCol()
	line := tz.r.NextLine()

	var s []rune // The full quoted string will be stored in this.

	q, err := tz.r.ReadChar()
	if err != nil {
		return nil, fmt.Errorf("Error reading quoted string.\n%s", err.Error())
	}
	tt, v := quoteTokenKind[q]
	if !v {
		return nil, fmt.Errorf("Invalid quote start char '%c'.", q)
	}
	s = append(s, q)

	done := false
	for true {
		c, err := tz.r.ReadChar()
		if err != nil {
			return nil, fmt.Errorf("Error reading quoted string.\n%s", err.Error())
		}

		switch c {
		case char.BackSlash:
			if raw {
				break
			}
			c, err = tz.esr.ReadChar(tz.r, tt)
			if err != nil {
				return nil, err
			}
		case char.NewLine:
			fallthrough
		case char.Return:
			if !raw {
				return nil, fmt.Errorf(
					"Unexpected newline while reading quoted string.")
			}
		case q:
			done = true
		default:
			// Do nothing
		}

		s = append(s, c)
		if done {
			break
		}
	}

	t := newToken(tt, s, line, col)
	return t, nil
}

func (tz *Tokenizer) readPyMultilineString() (*Token, error) {
	// Save the starting line and column for reporting.
	line := tz.r.NextLine()
	col := tz.r.NextCol()

	ss, err := tz.r.ReadSlice(3)
	if err != nil {
		return nil, fmt.Errorf("Error reading multine line string.\n%s", err.Error())
	}
	if string(ss) != TripleQuote {
		return nil, fmt.Errorf("Expecting '\"\"\"' as start of multiline string.")
	}

	var s []rune
	s = append(s, ss...)

	for true {
		c, err := tz.r.ReadChar()
		if err != nil {
			return nil, fmt.Errorf("Error reading multiline string.\n%s", err.Error())
		}

		s = append(s, c)

		if l := len(s); string(s[l-3:l]) == TripleQuote {
			break
		}
	}

	return newToken(token.PyMultilineString, s, line, col), nil
}
