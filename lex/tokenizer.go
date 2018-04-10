package lex

import (
	"fmt"
	"io"
	"uno/lex/char"
	"uno/lex/token"
)

var delimiterSet = map[rune]bool{
	'=':  true,
	'+':  true,
	'!':  true,
	'#':  true,
	'%':  true,
	'^':  true,
	'&':  true,
	'*':  true,
	'(':  true,
	')':  true,
	'-':  true,
	'|':  true,
	'\\': true,
	'\'': true,
	'"':  true,
	',':  true,
	'.':  true,
	':':  true,
	';':  true,
	'<':  true,
	'>':  true,
	'[':  true,
	']':  true,
	'{':  true,
	'}':  true,
}

func IsDelimiter(c rune) bool {
	_, e := delimiterSet[c]
	return e
}

func IsKeyword(s string) bool {
	_, e := KeywordMap[s]
	return e
}

func IsOperator(s string) bool {
	_, e := OperatorMap[s]
	return e
}

type Tokenizer struct {
	ts  token.KindSet
	r   *CharReader
	esr EscSeqReader

	// Convenience variables
	indent  bool // true if token.Indent is present in |ts|.
	newLine bool // true if token.NewLine is present in |ts|.
	tab     bool // true if token.Tab is present in |ts|.
}

// Returns the line on which the last successfully read or attempted
// character was present on. A value of 0 is returned before the first
// character is read.
func (tz *Tokenizer) Line() uint32 {
	return tz.r.Line()
}

// Returns the column on which the last successfully read or attempted
// character was present at. A value of 0 is returned before the first
// character is read.
func (tz *Tokenizer) Col() uint32 {
	return tz.r.Col()
}

// Returns a new Tokenizer object.
func NewTokenizer(r io.RuneReader, s token.KindSet, esr EscSeqReader) (*Tokenizer, error) {
	if r == nil {
		return nil, fmt.Errorf("A non-nil rune param is required.")
	}
	if s == nil {
		return nil, fmt.Errorf("A non-nil token.KindSet param is required.")
	}
	// Perform a sanity check of the input set of tokens.
	if s.Contains(token.Indent) && s.Contains(token.Tab) {
		return nil, fmt.Errorf("Tab and indent cannot be tokens together.")
	}
	if s.Contains(token.PySingleLineComment) && s.Contains(token.CPPDirective) {
		e := fmt.Errorf(
			"Python comments and C pre-processor directives cannot be tokens together.")
		return nil, e
	}
	if s.Contains(token.SingleQuoteString) && s.Contains(token.SingleQuoteCharacter) {
		e := fmt.Errorf(
			"Single quoted strings and character literals cannot be tokens together.")
		return nil, e
	}

	if esr == nil && (s.Contains(token.SingleQuoteString) || s.Contains(token.DoubleQuoteString)) {
		return nil, fmt.Errorf("A non-nil Escape Sequence Reader is required.")
	}

	tz := new(Tokenizer)
	tz.ts = s
	tz.r = NewCharReader(r)
	tz.esr = esr

	if s.Contains(token.Indent) {
		tz.indent = true
	}

	if s.Contains(token.Tab) {
		tz.tab = true
	}

	if s.Contains(token.NewLine) {
		tz.newLine = true
	}

	return tz, nil
}

// Returns true if there are further tokens, false otherwise.
func (tz *Tokenizer) HasNext() bool {
	_, e := tz.r.PeekChar()
	return e != io.EOF
}

// Returns the next token in the input.
// If an error occurs, it is not guaranteed to be recoverable.
func (tz *Tokenizer) NextToken() (*Token, error) {
	c, err := tz.r.PeekChar()
	if err != nil {
		return nil, err
	}

	switch {
	case c == char.Space:
		if tz.r.PreviousWasNewLine() && tz.indent {
			return tz.readIndentToken()
		}

		err = tz.skipSpace()
		if err != nil {
			return nil, err
		}

		return tz.NextToken()
	case c == char.Tab:
		if tz.r.PreviousWasNewLine() && tz.indent {
			return tz.readIndentToken()
		}

		if tz.tab {
			c, err = tz.r.ReadChar()
			if err != nil {
				err = fmt.Errorf("Error reading tab character.\n%s", err.Error())
				return nil, err
			}

			t := newToken(token.Tab, []rune{c}, tz.r.Line(), tz.r.Col())
			return t, nil
		}

		err = tz.skipSpace()
		if err != nil {
			return nil, err
		}

		return tz.NextToken()
	case c == char.NewLine || c == char.Return:
		c, err = tz.r.ReadChar()
		if err != nil {
			err = fmt.Errorf("Error reading new line character.\n%s", err.Error())
			return nil, err
		}

		if tz.newLine {
			t := newToken(
				token.NewLine, []rune{c}, tz.r.Line(), tz.r.Col())
			return t, nil
		}

		return tz.NextToken()
	case c == char.DoubleQuote:
		// It can either be the beginning of a double quoted string
		// or Python mutiline/doc string.
		if tz.ts.Contains(token.PyMultilineString) {
			q, err := tz.r.PeekSlice(3)
			if err == nil {
				if string(q) == TripleQuote {
					// It is a Python mutiline string.
					return tz.readPyMultilineString()
				}
			}
		}
		if tz.ts.Contains(token.DoubleQuoteString) {
			return tz.readQuotedString(false)
		}
	case c == char.SingleQuote:
		// It can either be a single quoted string or a character
		// literal.
		if tz.ts.Contains(token.SingleQuoteString) {
			return tz.readQuotedString(false)
		} else if tz.ts.Contains(token.SingleQuoteCharacter) {
			return tz.readSingleQuoteCharacter()
		}
	case c == char.BackQuote:
		if tz.ts.Contains(token.BackQuoteString) {
			return tz.readQuotedString(true)
		}
	case c == char.Hash:
		// It can either be a C pre-processor directive or a Python-style comment.
		if tz.ts.Contains(token.PySingleLineComment) {
			return tz.readPythonStyleComment()
		}
		if !tz.ts.Contains(token.CPPDirective) {
			break
		}

		cc, err := tz.r.PeekSlice(2)
		if err != nil || !isIdentifierBeginChar(cc[1]) {
			break
		}

		line := tz.r.NextLine()
		col := tz.r.NextCol()
		hash, err := tz.r.ReadChar()
		if err != nil {
			return nil, err
		}

		id, err := tz.readIdentifierString()
		if err != nil {
			err = fmt.Errorf(
				"Error reading preprocessor directive.\n%s",
				err.Error())
			return nil, err
		}

		s := []rune{hash}
		s = append(s, id...)
		t := newToken(token.CPPDirective, s, line, col)
		return t, nil
	case c == char.At:
		// Python style decorator.
		if !tz.ts.Contains(token.PythonDecorator) {
			break
		}

		cc, err := tz.r.PeekSlice(2)
		if err != nil || !isIdentifierBeginChar(cc[1]) {
			break
		}

		at, err := tz.r.ReadChar()
		if err != nil {
			return nil, err
		}

		line := tz.r.NextLine()
		col := tz.r.NextCol()
		id, err := tz.readIdentifierString()
		if err != nil {
			return nil, fmt.Errorf("Error reading Python decorator.\n%s", err.Error())
		}

		s := []rune{at}
		s = append(s, id...)
		return newToken(token.PythonDecorator, s, line, col), nil
	case c == char.Div:
		// It can either be the div operator itself or can be the C-style single
		// line comment or C-style multiline comment.
		cc, err := tz.r.PeekSlice(2)
		if err != nil {
			return tz.readOperator()
		}

		if cc[1] == char.Div && tz.ts.Contains(token.CSingleLineComment) {
			return tz.readCStyleSingleLineComment()
		} else if cc[1] == char.Mul && tz.ts.Contains(token.CMultiLineComment) {
			return tz.readCStyleMultiLineComment()
		} else {
			return tz.readOperator()
		}
	case c == char.Dot:
		// This can be the dot operator, or if it is followed by a number, then
		// a floating point number.
		cc, err := tz.r.PeekSlice(2)
		f := isDecimalDigit(cc[1]) || cc[1] == 'E' || cc[1] == 'e'
		if err == nil && f && tz.ts.Contains(token.FloatNumber) {
			return tz.readNumber()
		}
		return tz.readOperator()
	case isIdentifierBeginChar(c):
		return tz.readIdentifier()
	case isDecimalDigit(c):
		return tz.readNumber()
	case IsOperator(string(c)):
		return tz.readOperator()
	default:
		return nil, unExpectedCharacterError(c)
	}

	// If none of the above cases returned a token or an error, it means
	// that |c| is an unexpected character.
	return nil, unExpectedCharacterError(c)
}

func (tz *Tokenizer) newValidToken(t token.Kind, s []rune, l uint32, c uint32) (*Token, error) {
	if !tz.ts.Contains(t) {
		return nil, fmt.Errorf("Unexpected '%s'.", string(s))
	}

	return newToken(t, s, l, c), nil
}

func unExpectedCharacterError(c rune) error {
	return fmt.Errorf("Unexpected character '%c'.", c)
}
